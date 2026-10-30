#!/usr/bin/env python3
"""Compare isolated PostgreSQL/SQLite migration catalogs with checked-in source inventories."""

from __future__ import annotations

import argparse
import csv
import json
import re
import sys
from collections import defaultdict
from pathlib import Path


def read_csv(path: Path) -> list[dict[str, str]]:
    with path.open(newline="", encoding="utf-8") as stream:
        return list(csv.DictReader(stream))


def snake_case(value: str) -> str:
    initialisms = (
        "JWKS", "UUID", "HTTPS", "HTTP", "JSON", "URL", "URI", "ID",
        "API", "LLM", "MCP", "OCR", "TTS", "STT", "S3", "IP", "DNS",
        "JWT", "OIDC", "OAuth", "UI", "PDF", "CSV", "HTML", "XML",
        "SQL", "SSH", "TLS",
    )
    for token in sorted(initialisms, key=len, reverse=True):
        value = re.sub(token, token[0] + token[1:].lower(), value)
    first = re.sub(r"([A-Z]+)([A-Z][a-z])", r"\1_\2", value)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", first).lower()


def field_column(row: dict[str, str]) -> str | None:
    field = row["field"]
    tag = row.get("gorm_tag", "")
    if field in {"BaseModel", "ControlPlaneModel"}:
        return None
    if "-:migration" in tag:
        return None
    if re.search(r"(?:^|;)-(?:;|$)", tag) or "-:all" in tag:
        return None
    if row.get("column_tag"):
        return row["column_tag"]
    match = re.search(r"(?:^|;)column:([^;]+)", tag)
    return match.group(1) if match else snake_case(field)


def normalized_sql(value: str) -> str:
    return re.sub(r"[\s\"\x60]", "", value.lower())


def table_leaf(value: str) -> str:
    return value.rsplit(".", 1)[-1].strip('"')


def write_report(path: Path, lines: list[str]) -> None:
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--catalog-dir", required=True, type=Path)
    args = parser.parse_args()

    root = Path(__file__).resolve().parents[3]
    catalog_dir = args.catalog_dir.resolve()
    pg_path = catalog_dir / "postgres-catalog.json"
    sqlite_path = catalog_dir / "sqlite-catalog.json"
    for required in (pg_path, sqlite_path):
        if not required.is_file():
            raise SystemExit(f"missing physical catalog snapshot: {required}")

    tables = read_csv(root / "docs/migration/06-DATABASE-TABLE-MAPPING.csv")
    fields = read_csv(root / "docs/migration/06B-DATABASE-FIELD-INVENTORY.csv")
    gorm_indexes = read_csv(root / "docs/migration/06C-DATABASE-GORM-INDEX-MEMBERS.csv")
    postgres_ddl = read_csv(root / "docs/migration/06D-POSTGRES-CUSTOM-DDL.csv")
    vectors = read_csv(root / "docs/migration/06E-VECTOR-STORAGE-INVENTORY.csv")
    pg = json.loads(pg_path.read_text(encoding="utf-8"))
    sqlite = json.loads(sqlite_path.read_text(encoding="utf-8"))

    errors: list[str] = []
    warnings: list[str] = []
    model_tables = {row["table_name"] for row in tables}
    vector_tables = {
        row["table_name"] for row in vectors if row["engine"].lower().startswith("sqlite")
    }

    pg_tables = set(pg.get("tables", []))
    sqlite_tables = set(sqlite.get("tables", []))
    sqlite_application_tables = {
        name for name in sqlite_tables
        if not any(name.startswith(vector + "_") for vector in vector_tables)
    }
    for label, expected, actual in (
        ("PostgreSQL model tables", model_tables, pg_tables),
        ("SQLite model tables", model_tables, sqlite_application_tables),
        ("SQLite sqlite-vec tables", vector_tables, sqlite_tables),
    ):
        missing = sorted(expected - actual)
        if missing:
            errors.append(f"{label} missing: {', '.join(missing)}")
        extra = sorted(actual - expected) if "model tables" in label else []
        if extra:
            warnings.append(f"{label} extra: {', '.join(extra)}")

    field_columns: dict[tuple[str, str], str] = {}
    expected_pg_columns: dict[str, set[str]] = defaultdict(set)
    expected_sqlite_columns: dict[str, set[str]] = defaultdict(set)
    for row in fields:
        column = field_column(row)
        if not column:
            continue
        table = row["table_name"]
        field_columns[(table, row["field"])] = column
        expected_pg_columns[table].add(column)
        if "-:migration" not in row.get("gorm_tag", ""):
            expected_sqlite_columns[table].add(column)

    # PostgreSQL-only custom ALTER TABLE ADD COLUMN statements supplement the model fields.
    for row in postgres_ddl:
        if row["statement_type"] != "ALTER TABLE":
            continue
        match = re.search(
            r'ALTER\s+TABLE\s+"([^"]+)"\s+ADD\s+COLUMN(?:\s+IF\s+NOT\s+EXISTS)?\s+"([^"]+)"',
            row.get("sql", ""),
            re.IGNORECASE,
        )
        if match:
            expected_pg_columns[match.group(1)].add(match.group(2))

    pg_columns: dict[str, dict[str, dict[str, str]]] = defaultdict(dict)
    for column in pg.get("columns", []):
        pg_columns[column["table_name"]][column["column_name"]] = column
    sqlite_columns: dict[str, dict[str, dict[str, object]]] = defaultdict(dict)
    for column in sqlite.get("columns", []):
        sqlite_columns[column["table_name"]][column["name"]] = column

    def compare_columns(label: str, expected: dict[str, set[str]], actual: dict) -> None:
        missing: list[str] = []
        extra: list[str] = []
        for table in sorted(model_tables):
            expected_set = expected.get(table, set())
            actual_set = set(actual.get(table, {}))
            missing.extend(f"{table}.{column}" for column in sorted(expected_set - actual_set))
            extra.extend(f"{table}.{column}" for column in sorted(actual_set - expected_set))
        if missing:
            errors.append(f"{label} missing ({len(missing)}): {', '.join(missing[:30])}")
        if extra:
            warnings.append(f"{label} extra ({len(extra)}): {', '.join(extra[:30])}")

    compare_columns("PostgreSQL columns", expected_pg_columns, pg_columns)
    compare_columns("SQLite columns", expected_sqlite_columns, sqlite_columns)

    # Verify explicit NOT NULL and DEFAULT contracts from the GORM field inventory.
    nullability_errors: list[str] = []
    default_errors: list[str] = []
    for row in fields:
        column = field_column(row)
        if not column:
            continue
        table = row["table_name"]
        tag = row.get("gorm_tag", "")
        if "-:migration" in tag:
            sqlite_column = None
        else:
            sqlite_column = sqlite_columns.get(table, {}).get(column)
        pg_column = pg_columns.get(table, {}).get(column)
        if re.search(r"(?:^|;)not null(?:;|$)", tag) or "primaryKey" in tag:
            if pg_column and pg_column.get("is_nullable") != "NO":
                nullability_errors.append(f"PostgreSQL {table}.{column} is nullable")
            if sqlite_column and not (sqlite_column.get("not_null") or sqlite_column.get("primary_key")):
                nullability_errors.append(f"SQLite {table}.{column} is nullable")
        default_match = re.search(r"(?:^|;)default:([^;]+)", tag)
        if default_match and default_match.group(1).strip() not in {"", "(-)"}:
            if pg_column and not pg_column.get("column_default"):
                default_errors.append(f"PostgreSQL {table}.{column} has no default")
            if sqlite_column and not sqlite_column.get("default_value"):
                default_errors.append(f"SQLite {table}.{column} has no default")
    if nullability_errors:
        errors.extend(nullability_errors[:30])
        if len(nullability_errors) > 30:
            errors.append(f"... {len(nullability_errors) - 30} more nullability mismatches")
    if default_errors:
        errors.extend(default_errors[:30])
        if len(default_errors) > 30:
            errors.append(f"... {len(default_errors) - 30} more default mismatches")

    pg_index_rows = pg.get("indexes", [])
    sqlite_index_rows = sqlite.get("indexes", [])
    pg_index_map = {
        (row["table_name"], row["index_name"]): row["index_def"] for row in pg_index_rows
    }
    sqlite_index_map = {
        (row["table_name"], row["name"]): row for row in sqlite_index_rows
    }

    gorm_groups: dict[tuple[str, str], list[dict[str, str]]] = defaultdict(list)
    for row in gorm_indexes:
        if row["index_kind"] == "primaryKey" or row["index_name"] in {"", "GORM_DEFAULT"}:
            continue
        gorm_groups[(row["table_name"], row["index_name"])].append(row)
    # GORM's unnamed index tag uses idx_<table>_<column>; materialize those names
    # so the embedded BaseModel soft-delete indexes are verified as well.
    for row in fields:
        tag = row.get("gorm_tag", "")
        if not re.search(r"(?:^|;)index(?:;|$)", tag):
            continue
        if re.search(r"(?:^|;)(?:index|uniqueIndex):", tag):
            continue
        column = field_column(row)
        if not column:
            continue
        table = row["table_name"]
        name = f"idx_{table}_{column}"
        gorm_groups.setdefault((table, name), []).append({**row, "raw_tag": tag})

    missing_pg_indexes: list[str] = []
    missing_sqlite_indexes: list[str] = []
    index_definition_errors: list[str] = []
    for (table, name), members in sorted(gorm_groups.items()):
        pg_definition = pg_index_map.get((table, name))
        sqlite_definition = sqlite_index_map.get((table, name))
        if not pg_definition:
            missing_pg_indexes.append(f"{table}.{name}")
        if not sqlite_definition:
            missing_sqlite_indexes.append(f"{table}.{name}")
        members = sorted(
            enumerate(members),
            key=lambda pair: (int(pair[1]["priority"]) if pair[1].get("priority", "").isdigit() else 10_000, pair[0]),
        )
        expected_columns = [
            field_columns.get((table, member["field"]), snake_case(member["field"]))
            for _, member in members
        ]
        if pg_definition and len(expected_columns) > 1:
            normalized = normalized_sql(pg_definition)
            needle = "(" + ",".join(expected_columns) + ")"
            if needle not in normalized:
                index_definition_errors.append(f"PostgreSQL {table}.{name} column order differs from {expected_columns}")
        if sqlite_definition and len(expected_columns) > 1:
            actual_columns = sqlite_definition.get("columns", [])
            if actual_columns != expected_columns:
                index_definition_errors.append(
                    f"SQLite {table}.{name} columns {actual_columns} differ from {expected_columns}"
                )
        if any("where:" in member.get("raw_tag", member.get("gorm_tag", "")).lower() for _, member in members):
            if pg_definition and "where" not in pg_definition.lower():
                index_definition_errors.append(f"PostgreSQL partial index {table}.{name} has no predicate")
            if sqlite_definition and "where" not in sqlite_definition.get("sql", "").lower():
                index_definition_errors.append(f"SQLite partial index {table}.{name} has no predicate")

    if missing_pg_indexes:
        errors.append(f"PostgreSQL GORM indexes missing ({len(missing_pg_indexes)}): {', '.join(missing_pg_indexes[:30])}")
    if missing_sqlite_indexes:
        errors.append(f"SQLite GORM indexes missing ({len(missing_sqlite_indexes)}): {', '.join(missing_sqlite_indexes[:30])}")
    errors.extend(index_definition_errors[:30])
    if len(index_definition_errors) > 30:
        errors.append(f"... {len(index_definition_errors) - 30} more index-definition mismatches")

    custom_indexes = [
        row for row in postgres_ddl
        if row["statement_type"] in {"CREATE INDEX", "CREATE UNIQUE INDEX"}
        and row.get("object_name")
        and row.get("table_name")
        and "%" not in row.get("sql", "")
        and "{table}" not in row.get("sql", "")
    ]
    custom_index_errors: list[str] = []
    for row in custom_indexes:
        key = (row["table_name"], row["object_name"])
        definition = pg_index_map.get(key)
        if not definition:
            custom_index_errors.append(f"PostgreSQL custom index missing: {key[0]}.{key[1]}")
            continue
        normalized = normalized_sql(definition)
        expected_columns = re.findall(r'"([^"]+)"', row.get("columns", ""))
        if expected_columns:
            using = re.search(r"\bUSING\s+\w+\s*\((.*)\)(?:\s+WHERE\s+.*)?$", definition, re.IGNORECASE)
            section = using.group(1) if using else definition[definition.find("(") + 1:definition.rfind(")")]
            position = 0
            for column in expected_columns:
                match = re.search(
                    rf"(?<![a-z0-9_]){re.escape(column.lower())}(?![a-z0-9_])",
                    section[position:],
                    re.IGNORECASE,
                )
                if not match:
                    custom_index_errors.append(f"PostgreSQL custom index {key[1]} columns differ from inventory")
                    break
                position += match.end()
        if row.get("predicate") and "where" not in definition.lower():
            custom_index_errors.append(f"PostgreSQL custom partial index {key[1]} has no predicate")
    errors.extend(custom_index_errors[:30])
    if len(custom_index_errors) > 30:
        errors.append(f"... {len(custom_index_errors) - 30} more custom-index mismatches")

    legacy_index_key = ("chat_conversation_projects", "idx_chat_conversation_projects_public_id")
    if legacy_index_key in pg_index_map:
        errors.append("PostgreSQL legacy full unique index remains; expected active-row partial index")
    if legacy_index_key in sqlite_index_map:
        errors.append("SQLite legacy full unique index remains; expected active-row partial index")

    pg_pk_tables = {
        table_leaf(row["table_name"])
        for row in pg.get("constraints", [])
        if row.get("constraint_type") == "p"
    }
    sqlite_pk_tables = {
        table for table, columns in sqlite_columns.items()
        if any(column.get("primary_key", 0) for column in columns.values())
    }
    missing_pg_pk = sorted(model_tables - pg_pk_tables)
    missing_sqlite_pk = sorted(model_tables - sqlite_pk_tables)
    if missing_pg_pk:
        errors.append(f"PostgreSQL primary keys missing: {', '.join(missing_pg_pk)}")
    if missing_sqlite_pk:
        errors.append(f"SQLite primary keys missing: {', '.join(missing_sqlite_pk)}")

    sqlite_objects = {row["name"]: row for row in sqlite.get("objects", []) if row.get("type") == "table"}
    vector_errors: list[str] = []
    for row in vectors:
        name = row["table_name"]
        obj = sqlite_objects.get(name)
        if not obj or "vec0" not in obj.get("sql", "").lower():
            vector_errors.append(f"SQLite vector virtual table missing or not vec0: {name}")
            continue
        ddl = obj.get("sql", "").lower()
        required_tokens = [
            *[part.strip().lower() for part in row["key_columns"].split(",") if part.strip()],
            row["partition_column"].lower(),
            row["signature_column"].lower(),
            f"float[{row['embedding_dimensions']}]",
            f"distance_metric={row['distance_metric'].lower()}",
        ]
        for token in required_tokens:
            if token and token not in ddl:
                vector_errors.append(f"SQLite vector table {name} missing DDL token {token}")
    errors.extend(vector_errors)

    report = [
        "# Isolated database migration catalog comparison",
        "",
        f"- PostgreSQL schema: `{pg.get('schema', '')}`",
        f"- PostgreSQL model tables: {len(pg_tables)} (expected {len(model_tables)})",
        f"- PostgreSQL columns: {len(pg.get('columns', []))}",
        f"- PostgreSQL indexes: {len(pg_index_rows)}",
        f"- PostgreSQL constraints: {len(pg.get('constraints', []))}",
        f"- SQLite tables: {len(sqlite_tables)} physical including sqlite-vec shadow tables (expected {len(model_tables) + len(vector_tables)} application/virtual tables)",
        f"- SQLite columns: {len(sqlite.get('columns', []))}",
        f"- SQLite indexes: {len(sqlite_index_rows)}",
        f"- GORM named indexes checked: {len(gorm_groups)}",
        f"- PostgreSQL custom indexes checked: {len(custom_indexes)}",
        "",
        f"## Result: {'FAIL' if errors else 'PASS'}",
        "",
    ]
    if warnings:
        report += ["## Warnings", "", *[f"- {item}" for item in warnings], ""]
    if errors:
        report += ["## Mismatches", "", *[f"- {item}" for item in errors], ""]
    else:
        report += [
            "All model tables and inventoried columns are present in both engines.",
            "GORM and PostgreSQL custom indexes, partial-index predicates, primary keys, declared nullability/defaults and sqlite-vec virtual tables match the checked-in source inventories.",
            "",
        ]
    report_path = catalog_dir / "catalog-diff-report.md"
    write_report(report_path, report)
    summary_path = Path(__import__("os").environ.get("GITHUB_STEP_SUMMARY", ""))
    if str(summary_path) not in {"", "."}:
        with summary_path.open("a", encoding="utf-8") as stream:
            stream.write("\n".join(report) + "\n")
    print("\n".join(report))
    if errors:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())

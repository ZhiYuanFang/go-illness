#!/usr/bin/env python3
"""合并 openspec/specs 下全部 capability spec.md 为单一版本基线文件。

无 capability 目录时写出合法空壳合并文件（与 bootstrap v0.0.1 形态一致）。
"""
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SPECS_DIR = ROOT / "openspec" / "specs"
DEFAULT_VERSION = "v0.0.1"


def write_empty_shell(out: Path, version: str) -> None:
    """写出无 capability 的空壳基线。"""
    parts: list[str] = []
    parts.append(f"# OpenSpec 基线规格 {version}\n\n")
    parts.append(
        f"> 本文件由 `openspec/specs` 下全部 capability 规格于 **{version}** 合并而成，"
        "作为该版本的确定性规则基线，便于按版本查阅。\n\n"
    )
    parts.append(
        "> 后续新变更请基于本文件对照增量，或通过 OpenSpec 新建 change 再合并至下一版本规格。\n\n"
    )
    parts.append("## 目录\n\n")
    parts.append(">（本版为空壳：尚无 capability。）\n\n")
    parts.append("---\n\n")
    parts.append(
        f"**本版说明**：`{version}` 为空壳基线，不含任何 capability 章节。"
        "业务与工程能力经 change 增量引入后，通过 `/opsx-archive vX.Y.Z` 收束至后续版本。\n"
    )
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text("".join(parts), encoding="utf-8", newline="\n")


def main() -> None:
    import sys

    version = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_VERSION
    out = SPECS_DIR / version / "spec.md"

    files = sorted(SPECS_DIR.glob("*/spec.md"))
    # 排除已生成的版本合并基线目录（如 v0.0.1、v1.0.3）
    files = [f for f in files if not re.match(r"^v\d+\.\d+\.\d+$", f.parent.name)]
    if not files:
        write_empty_shell(out, version)
        text = out.read_text(encoding="utf-8")
        print(f"merged 0 specs (empty shell) -> {out.relative_to(ROOT)}")
        print(f"lines={len(text.splitlines())} bytes={out.stat().st_size}")
        return

    parts: list[str] = []
    parts.append(f"# OpenSpec 基线规格 {version}\n\n")
    parts.append(
        f"> 本文件由 `openspec/specs` 下全部 capability 规格于 **{version}** 合并而成，"
        "作为该版本的确定性规则基线，便于按版本查阅。\n\n"
    )
    parts.append(
        "> 后续新变更请基于本文件对照增量，或通过 OpenSpec 新建 change 再合并至下一版本规格。\n\n"
    )
    parts.append("## 目录\n\n")
    for f in files:
        cap = f.parent.name
        parts.append(f"- [{cap}](#{cap})\n")
    parts.append("\n---\n\n")

    for f in files:
        cap = f.parent.name
        content = f.read_text(encoding="utf-8").strip()
        parts.append(f"## {cap}\n\n")
        parts.append(f"<!-- source: openspec/specs/{cap}/spec.md -->\n\n")
        parts.append(content)
        parts.append("\n\n---\n\n")

    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text("".join(parts), encoding="utf-8", newline="\n")
    text = out.read_text(encoding="utf-8")
    print(f"merged {len(files)} specs -> {out.relative_to(ROOT)}")
    print(f"lines={len(text.splitlines())} bytes={out.stat().st_size}")


if __name__ == "__main__":
    main()

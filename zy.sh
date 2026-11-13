#!/usr/bin/env bash
set -euo pipefail

die(){ echo "❌ $*" >&2; exit 1; }
warn(){ echo "⚠️  $*" >&2; }

git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "当前目录不是 git 仓库"

read -r -p "开始时间 since（如 2026-02-01 或 '2 weeks ago'）： " SINCE
read -r -p "结束时间 until（如 2026-02-28 或 'today'）： " UNTIL

[[ -n "${SINCE:-}" ]] || die "since 不能为空"
[[ -n "${UNTIL:-}" ]] || die "until 不能为空"

echo
echo "输入作者（author）匹配条件：可以是名字或邮箱；支持正则。"
echo "建议用邮箱更精准。"
echo "一行一个，输入空行结束："

AUTHORS=()
while true; do
  read -r -p "author> " a || true
  [[ -z "${a:-}" ]] && break
  AUTHORS+=("$a")
done

[[ ${#AUTHORS[@]} -gt 0 ]] || die "至少输入 1 个 author"

read -r -p "可选：只过滤某个目录前缀（如 src/，回车跳过）： " PATH_PREFIX || true
read -r -p "可选：只过滤某个后缀（如 .go / .java，回车跳过）： " SUFFIX || true

echo
read -r -p "输出是否排序（默认 Y）[Y/n]： " SORT_OPT || true
SORT_OPT="${SORT_OPT:-Y}"

# author1|author2|...
AUTHOR_REGEX="$(printf "%s|" "${AUTHORS[@]}")"
AUTHOR_REGEX="${AUTHOR_REGEX%|}"

# 获取文件列表（name-only），去掉空行
OUT="$(
  git log \
    --since="$SINCE" \
    --until="$UNTIL" \
    --author="$AUTHOR_REGEX" \
    --name-only \
    --pretty=format: \
  | sed '/^$/d'
)"

# 目录过滤（前缀）
if [[ -n "${PATH_PREFIX:-}" ]]; then
  OUT="$(printf "%s\n" "$OUT" | awk -v p="$PATH_PREFIX" 'index($0,p)==1')"
fi

# 后缀过滤
if [[ -n "${SUFFIX:-}" ]]; then
  # 末尾匹配
  OUT="$(printf "%s\n" "$OUT" | awk -v s="$SUFFIX" 'substr($0,length($0)-length(s)+1)==s')"
fi

# 去重（可选排序）
if [[ "$SORT_OPT" =~ ^([Yy]|yes|YES)$ ]]; then
  OUT="$(printf "%s\n" "$OUT" | sort -u)"
else
  # 不排序的去重：保留首次出现顺序（awk 去重）
  OUT="$(printf "%s\n" "$OUT" | awk '!seen[$0]++')"
fi

echo
echo "====== 结果（去重后的文件路径）======"
if [[ -n "${OUT:-}" ]]; then
  printf "%s\n" "$OUT"
  echo
  echo "总数：$(printf "%s\n" "$OUT" | wc -l | tr -d ' ')"
else
  echo "(空) 没有匹配到文件。"
fi

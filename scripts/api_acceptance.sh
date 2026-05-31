#!/usr/bin/env bash
set -euo pipefail

API_BASE="${API_BASE:-http://127.0.0.1:8091/api}"
DIST_EMAIL="${DIST_EMAIL:-dist_demo@example.com}"
DIST_PASSWORD="${DIST_PASSWORD:-Viewer123!}"
OPS_EMAIL="${OPS_EMAIL:-operator_demo@local.dev}"
OPS_PASSWORD="${OPS_PASSWORD:-Distributor123!}"

WORKDIR="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="$WORKDIR/test-results"
mkdir -p "$OUT_DIR"
REPORT="$OUT_DIR/api-acceptance-$(date +%Y%m%d-%H%M%S).md"
LATEST_REPORT="$OUT_DIR/api-acceptance-latest.md"

LAST_BODY=""
LAST_STATUS=""

request() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local token="${4:-}"
  local headers=(-H "Content-Type: application/json")
  if [[ -n "$token" ]]; then
    headers+=(-H "Authorization: Bearer $token")
  fi
  if [[ -n "$body" ]]; then
    LAST_BODY="$(curl -sS -X "$method" "${headers[@]}" "$API_BASE$path" -d "$body" -w $'\n%{http_code}')"
  else
    LAST_BODY="$(curl -sS -X "$method" "${headers[@]}" "$API_BASE$path" -w $'\n%{http_code}')"
  fi
  LAST_STATUS="$(printf '%s' "$LAST_BODY" | tail -n1)"
  LAST_BODY="$(printf '%s' "$LAST_BODY" | sed '$d')"
}

assert_status() {
  local expected="$1"
  if [[ "$LAST_STATUS" != "$expected" ]]; then
    echo "Expected status $expected, got $LAST_STATUS" >&2
    echo "$LAST_BODY" >&2
    exit 1
  fi
}

json_get() {
  local expr="$1"
  printf '%s' "$LAST_BODY" | jq -r "$expr"
}

append_report() {
  local title="$1"
  local method="$2"
  local path="$3"
  {
    echo "## $title"
    echo
    echo "- Method: \`$method\`"
    echo "- Path: \`$path\`"
    echo "- Status: \`$LAST_STATUS\`"
    echo
    echo '```json'
    printf '%s\n' "$LAST_BODY"
    echo '```'
    echo
  } >>"$REPORT"
}

echo "# sub2api-distributor API 验收记录" >"$REPORT"
echo >>"$REPORT"
echo "- API Base: \`$API_BASE\`" >>"$REPORT"
echo "- 执行时间: \`$(date '+%Y-%m-%d %H:%M:%S %z')\`" >>"$REPORT"
echo >>"$REPORT"

request POST /auth/login "{\"email\":\"$DIST_EMAIL\",\"password\":\"$DIST_PASSWORD\"}"
assert_status 200
DIST_TOKEN="$(json_get '.token')"
append_report "分销商登录" POST /auth/login

request GET /me "" "$DIST_TOKEN"
assert_status 200
DIST_USER_ID="$(json_get '.user_id')"
append_report "分销商 me" GET /me

request GET /portal/dashboard "" "$DIST_TOKEN"
assert_status 200
append_report "分销商看板" GET /portal/dashboard

request GET /portal/invite-meta "" "$DIST_TOKEN"
assert_status 200
append_report "分销商邀请码信息" GET /portal/invite-meta

request GET /portal/invitees "" "$DIST_TOKEN"
assert_status 200
append_report "分销商邀请用户列表" GET /portal/invitees

request GET /portal/rebates "" "$DIST_TOKEN"
assert_status 200
append_report "分销商返利明细" GET /portal/rebates

request GET /portal/withdrawals "" "$DIST_TOKEN"
assert_status 200
PAYING_WITHDRAW_ID="$(json_get '.items[] | select(.status=="paying") | .id' | head -n1)"
append_report "分销商提现列表" GET /portal/withdrawals

request GET /portal/settlement-profile "" "$DIST_TOKEN"
assert_status 200
append_report "分销商收款信息查询" GET /portal/settlement-profile

request PUT /portal/settlement-profile '{"display_name":"演示分销商 A","status":"active","settlement_channel":"alipay","settlement_account_name":"邀请人A","settlement_account_no":"viewer-inviter@demo","settlement_account_extra":"{\"nickname\":\"邀请人A\",\"updated_by\":\"acceptance\"}","notes":"acceptance update"}' "$DIST_TOKEN"
assert_status 200
append_report "分销商收款信息更新" PUT /portal/settlement-profile

request POST /portal/withdrawals '{"amount":1.23,"remark":"acceptance create"}' "$DIST_TOKEN"
assert_status 200
NEW_WITHDRAW_ID="$(json_get '.id')"
append_report "分销商发起提现" POST /portal/withdrawals

request POST "/portal/withdrawals/$NEW_WITHDRAW_ID/cancel" "" "$DIST_TOKEN"
assert_status 200
append_report "分销商取消自己的提现" POST "/portal/withdrawals/$NEW_WITHDRAW_ID/cancel"

request POST /auth/logout "" "$DIST_TOKEN"
assert_status 200
append_report "分销商登出" POST /auth/logout

request POST /auth/login "{\"email\":\"$OPS_EMAIL\",\"password\":\"$OPS_PASSWORD\"}"
assert_status 200
OPS_TOKEN="$(json_get '.token')"
append_report "运营登录" POST /auth/login

request GET /me "" "$OPS_TOKEN"
assert_status 200
append_report "运营 me" GET /me

request GET /ops/distributors "" "$OPS_TOKEN"
assert_status 200
append_report "运营分销商列表" GET /ops/distributors

request GET "/ops/users/lookup?q=dist_demo" "" "$OPS_TOKEN"
assert_status 200
append_report "运营主系统用户搜索" GET "/ops/users/lookup?q=dist_demo"

request GET "/ops/distributors/$DIST_USER_ID" "" "$OPS_TOKEN"
assert_status 200
append_report "运营查看分销商详情" GET "/ops/distributors/$DIST_USER_ID"

request PUT "/ops/distributors/$DIST_USER_ID/profile" '{"display_name":"演示分销商 A","status":"active","settlement_channel":"alipay","settlement_account_name":"邀请人A","settlement_account_no":"viewer-inviter@demo","settlement_account_extra":"{\"nickname\":\"邀请人A\",\"updated_by\":\"ops\"}","notes":"ops update"}' "$OPS_TOKEN"
assert_status 200
append_report "运营更新分销商资料" PUT "/ops/distributors/$DIST_USER_ID/profile"

request GET /ops/withdrawals "" "$OPS_TOKEN"
assert_status 200
OPS_PAYING_ID="$(json_get '.items[] | select(.status=="paying") | .id' | head -n1)"
append_report "运营提现列表" GET /ops/withdrawals

if [[ -z "$OPS_PAYING_ID" ]]; then
  echo "No paying withdrawal found for operator flow" >&2
  exit 1
fi

request GET "/ops/withdrawals/$OPS_PAYING_ID" "" "$OPS_TOKEN"
assert_status 200
append_report "运营提现详情" GET "/ops/withdrawals/$OPS_PAYING_ID"

request POST "/ops/withdrawals/$OPS_PAYING_ID/mark-paid" '{"paid_channel":"manual","paid_reference_no":"OPS-ACCEPT-0001","paid_remark":"acceptance mark paid"}' "$OPS_TOKEN"
assert_status 200
append_report "运营标记已打款" POST "/ops/withdrawals/$OPS_PAYING_ID/mark-paid"

request POST /auth/login "{\"email\":\"$DIST_EMAIL\",\"password\":\"$DIST_PASSWORD\"}"
assert_status 200
DIST_TOKEN="$(json_get '.token')"
append_report "分销商重新登录" POST /auth/login

request POST /portal/withdrawals '{"amount":0.5,"remark":"acceptance ops cancel target"}' "$DIST_TOKEN"
assert_status 200
OPS_CANCEL_TARGET_ID="$(json_get '.id')"
append_report "创建运营取消目标提现单" POST /portal/withdrawals

request POST "/ops/withdrawals/$OPS_CANCEL_TARGET_ID/cancel" "" "$OPS_TOKEN"
assert_status 200
append_report "运营取消提现单" POST "/ops/withdrawals/$OPS_CANCEL_TARGET_ID/cancel"

request POST /auth/logout "" "$OPS_TOKEN"
assert_status 200
append_report "运营登出" POST /auth/logout

cp "$REPORT" "$LATEST_REPORT"

echo "Acceptance report written to: $REPORT"
echo "Latest acceptance report copied to: $LATEST_REPORT"

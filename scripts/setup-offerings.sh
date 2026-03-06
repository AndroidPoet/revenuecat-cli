#!/usr/bin/env bash
set -euo pipefail

# ──────────────────────────────────────────────────────────────────────────────
# setup-offerings.sh — One-click RevenueCat offerings setup
#
# Interactive menu-driven setup that creates the full RevenueCat stack:
# apps → products → entitlements → offerings → packages
#
# Features:
#   - Choose platform: iOS only, Android only, or both
#   - Add pricing info per product (display reference)
#   - Preset templates (freemium, paywall, metered, etc.)
#   - Config file support for repeatable setups
#   - Dry-run mode to preview without creating
#
# Usage:
#   ./scripts/setup-offerings.sh                     # Interactive menu
#   ./scripts/setup-offerings.sh --config plan.yaml  # From config file
#   ./scripts/setup-offerings.sh --dry-run            # Preview without creating
# ──────────────────────────────────────────────────────────────────────────────

RC="${RC_BIN:-rc}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# State
DRY_RUN=false
CONFIG_FILE=""
VERBOSE=false
PLATFORM=""  # "ios", "android", "both"

# Created resource IDs
IOS_APP_ID=""
ANDROID_APP_ID=""
OFFERING_ID=""
declare -a PACKAGE_IDS=()

# ── Helpers ──────────────────────────────────────────────────────────────────

info()    { echo -e "${BLUE}▸${NC} $*"; }
ok()      { echo -e "${GREEN}✓${NC} $*"; }
warn()    { echo -e "${YELLOW}!${NC} $*"; }
err()     { echo -e "${RED}✗${NC} $*" >&2; }
step()    { echo -e "\n${BOLD}${CYAN}[$1/$TOTAL_STEPS]${NC} ${BOLD}$2${NC}"; }
divider() { echo -e "${DIM}────────────────────────────────────────────────────────${NC}"; }

run_rc() {
  if $DRY_RUN; then
    echo -e "  ${DIM}(dry-run)${NC} $RC $*"
    return 0
  fi
  if $VERBOSE; then
    echo -e "  ${DIM}\$${NC} $RC $*"
  fi
  $RC "$@"
}

extract_id() {
  grep -oE '"id"\s*:\s*"[^"]+"' | head -1 | grep -oE '"[^"]+"\s*$' | tr -d '"' | xargs
}

prompt() {
  local var="$1" prompt_text="$2" default="${3:-}"
  if [[ -n "$default" ]]; then
    read -rp "$(echo -e "${CYAN}?${NC}") $prompt_text [${default}]: " input
    eval "$var=\"${input:-$default}\""
  else
    read -rp "$(echo -e "${CYAN}?${NC}") $prompt_text: " input
    eval "$var=\"$input\""
  fi
}

confirm() {
  read -rp "$(echo -e "${YELLOW}?${NC}") $1 (y/N): " yn
  [[ "$yn" =~ ^[Yy]$ ]]
}

pick_one() {
  local prompt_text="$1"
  shift
  local options=("$@")
  echo ""
  for i in "${!options[@]}"; do
    echo -e "  ${CYAN}$((i+1)))${NC} ${options[$i]}"
  done
  echo ""
  local choice
  read -rp "$(echo -e "${CYAN}?${NC}") $prompt_text: " choice
  echo "$choice"
}

banner() {
  echo ""
  echo -e "${BOLD}${CYAN}"
  echo "  ┌─────────────────────────────────────────────┐"
  echo "  │         RevenueCat One-Click Setup           │"
  echo "  │                                              │"
  echo "  │  Apps  Products  Entitlements  Offerings      │"
  echo "  │  ────  ────────  ────────────  ─────────      │"
  echo "  │  All wired together, ready to go.            │"
  echo "  └─────────────────────────────────────────────┘"
  echo -e "${NC}"
}

# ── Platform Selection ───────────────────────────────────────────────────────

select_platform() {
  echo -e "\n${BOLD}Platform${NC}"
  echo ""
  echo -e "  ${CYAN}1)${NC} iOS only           (App Store)"
  echo -e "  ${CYAN}2)${NC} Android only       (Play Store)"
  echo -e "  ${CYAN}3)${NC} Both iOS + Android  (recommended)"
  echo ""
  local choice
  read -rp "$(echo -e "${CYAN}?${NC}") Select platform: " choice
  case "$choice" in
    1) PLATFORM="ios" ;;
    2) PLATFORM="android" ;;
    3|"") PLATFORM="both" ;;
    *) PLATFORM="both" ;;
  esac
  ok "Platform: $PLATFORM"
}

# ── Preset Templates ────────────────────────────────────────────────────────

select_template() {
  echo -e "\n${BOLD}Setup Template${NC}"
  echo ""
  echo -e "  ${CYAN}1)${NC} ${BOLD}Freemium${NC}          Weekly, Monthly, Annual subscriptions + Lifetime"
  echo -e "  ${CYAN}2)${NC} ${BOLD}Simple Paywall${NC}    Monthly + Annual subscriptions"
  echo -e "  ${CYAN}3)${NC} ${BOLD}Trial-First${NC}       Free trial → Monthly, Annual (with trial periods)"
  echo -e "  ${CYAN}4)${NC} ${BOLD}Tiered${NC}            Basic + Pro tiers, each with Monthly + Annual"
  echo -e "  ${CYAN}5)${NC} ${BOLD}Consumable${NC}        Credits / coin packs (one-time purchases)"
  echo -e "  ${CYAN}6)${NC} ${BOLD}Custom${NC}            Build your own from scratch"
  echo ""
  local choice
  read -rp "$(echo -e "${CYAN}?${NC}") Choose template (1-6): " choice
  echo "$choice"
}

apply_template() {
  local template="$1" base_id="$2"

  case "$template" in
    1) # Freemium
      TEMPLATE_ENTITLEMENTS="premium"
      TEMPLATE_OFFERING_KEY="default"
      TEMPLATE_OFFERING_NAME="Default Offering"
      TEMPLATE_PKG_KEYS=("weekly" "monthly" "annual" "lifetime")
      TEMPLATE_PKG_NAMES=("Weekly" "Monthly" "Annual" "Lifetime")
      TEMPLATE_PKG_TYPES=("subscription" "subscription" "subscription" "one_time")
      TEMPLATE_PKG_PRICES=("\$2.99/week" "\$9.99/month" "\$49.99/year" "\$149.99")
      TEMPLATE_PKG_IOS_IDS=("${base_id}.weekly" "${base_id}.monthly" "${base_id}.annual" "${base_id}.lifetime")
      TEMPLATE_PKG_ANDROID_IDS=("${base_id}.weekly" "${base_id}.monthly" "${base_id}.annual" "${base_id}.lifetime")
      TEMPLATE_PKG_ENTS=("premium" "premium" "premium" "premium")
      ;;
    2) # Simple Paywall
      TEMPLATE_ENTITLEMENTS="premium"
      TEMPLATE_OFFERING_KEY="default"
      TEMPLATE_OFFERING_NAME="Default Offering"
      TEMPLATE_PKG_KEYS=("monthly" "annual")
      TEMPLATE_PKG_NAMES=("Monthly" "Annual")
      TEMPLATE_PKG_TYPES=("subscription" "subscription")
      TEMPLATE_PKG_PRICES=("\$9.99/month" "\$49.99/year")
      TEMPLATE_PKG_IOS_IDS=("${base_id}.monthly" "${base_id}.annual")
      TEMPLATE_PKG_ANDROID_IDS=("${base_id}.monthly" "${base_id}.annual")
      TEMPLATE_PKG_ENTS=("premium" "premium")
      ;;
    3) # Trial-First
      TEMPLATE_ENTITLEMENTS="premium"
      TEMPLATE_OFFERING_KEY="default"
      TEMPLATE_OFFERING_NAME="Default Offering"
      TEMPLATE_PKG_KEYS=("monthly" "annual")
      TEMPLATE_PKG_NAMES=("Monthly (7-day trial)" "Annual (14-day trial)")
      TEMPLATE_PKG_TYPES=("subscription" "subscription")
      TEMPLATE_PKG_PRICES=("\$9.99/month" "\$59.99/year")
      TEMPLATE_PKG_IOS_IDS=("${base_id}.monthly.trial" "${base_id}.annual.trial")
      TEMPLATE_PKG_ANDROID_IDS=("${base_id}.monthly.trial" "${base_id}.annual.trial")
      TEMPLATE_PKG_ENTS=("premium" "premium")
      ;;
    4) # Tiered
      TEMPLATE_ENTITLEMENTS="basic,pro"
      TEMPLATE_OFFERING_KEY="default"
      TEMPLATE_OFFERING_NAME="Default Offering"
      TEMPLATE_PKG_KEYS=("basic_monthly" "basic_annual" "pro_monthly" "pro_annual")
      TEMPLATE_PKG_NAMES=("Basic Monthly" "Basic Annual" "Pro Monthly" "Pro Annual")
      TEMPLATE_PKG_TYPES=("subscription" "subscription" "subscription" "subscription")
      TEMPLATE_PKG_PRICES=("\$4.99/month" "\$29.99/year" "\$14.99/month" "\$99.99/year")
      TEMPLATE_PKG_IOS_IDS=("${base_id}.basic.monthly" "${base_id}.basic.annual" "${base_id}.pro.monthly" "${base_id}.pro.annual")
      TEMPLATE_PKG_ANDROID_IDS=("${base_id}.basic.monthly" "${base_id}.basic.annual" "${base_id}.pro.monthly" "${base_id}.pro.annual")
      TEMPLATE_PKG_ENTS=("basic" "basic" "basic,pro" "basic,pro")
      ;;
    5) # Consumable
      TEMPLATE_ENTITLEMENTS="credits"
      TEMPLATE_OFFERING_KEY="default"
      TEMPLATE_OFFERING_NAME="Credit Packs"
      TEMPLATE_PKG_KEYS=("small_pack" "medium_pack" "large_pack" "mega_pack")
      TEMPLATE_PKG_NAMES=("10 Credits" "50 Credits" "150 Credits" "500 Credits")
      TEMPLATE_PKG_TYPES=("one_time" "one_time" "one_time" "one_time")
      TEMPLATE_PKG_PRICES=("\$0.99" "\$3.99" "\$9.99" "\$24.99")
      TEMPLATE_PKG_IOS_IDS=("${base_id}.credits.10" "${base_id}.credits.50" "${base_id}.credits.150" "${base_id}.credits.500")
      TEMPLATE_PKG_ANDROID_IDS=("${base_id}.credits.10" "${base_id}.credits.50" "${base_id}.credits.150" "${base_id}.credits.500")
      TEMPLATE_PKG_ENTS=("credits" "credits" "credits" "credits")
      ;;
  esac
}

# ── Config File Support ──────────────────────────────────────────────────────

generate_sample_config() {
  cat <<'YAML'
# setup-offerings.yaml — RevenueCat one-click setup config
#
# Usage: ./scripts/setup-offerings.sh --config setup-offerings.yaml

# Platform: "ios", "android", or "both"
platform: "both"

# App details
app_name: "My App"
ios_bundle_id: "com.example.myapp"
android_package_name: "com.example.myapp"

# Entitlements to create
entitlements:
  - lookup_key: premium
    display_name: "Premium"
  - lookup_key: pro
    display_name: "Pro"

# Offering
offering:
  lookup_key: default
  display_name: "Default Offering"
  set_current: true

# Packages and their products
packages:
  - lookup_key: weekly
    display_name: "Weekly"
    ios_product_id: "com.example.myapp.weekly"
    android_product_id: "com.example.myapp.weekly"
    type: subscription
    price: "$2.99/week"
    entitlements:
      - premium

  - lookup_key: monthly
    display_name: "Monthly"
    ios_product_id: "com.example.myapp.monthly"
    android_product_id: "com.example.myapp.monthly"
    type: subscription
    price: "$9.99/month"
    entitlements:
      - premium

  - lookup_key: annual
    display_name: "Annual"
    ios_product_id: "com.example.myapp.annual"
    android_product_id: "com.example.myapp.annual"
    type: subscription
    price: "$49.99/year"
    entitlements:
      - premium
      - pro

  - lookup_key: lifetime
    display_name: "Lifetime"
    ios_product_id: "com.example.myapp.lifetime"
    android_product_id: "com.example.myapp.lifetime"
    type: one_time
    price: "$149.99"
    entitlements:
      - premium
      - pro
YAML
}

yaml_val() {
  local key="$1" file="$2"
  grep "^${key}:" "$file" | head -1 | sed 's/^[^:]*:\s*//' | sed 's/^"\(.*\)"$/\1/' | xargs
}

yaml_block_count() {
  local key="$1" file="$2"
  awk "/^${key}:/{found=1; next} found && /^  - /{count++} found && /^[^ ]/{exit} END{print count+0}" "$file"
}

yaml_block_field() {
  local block_key="$1" index="$2" field="$3" file="$4"
  awk -v idx="$index" -v fld="$field" '
    /^'"$block_key"':/ { found=1; next }
    found && /^  - / { count++ }
    found && count==idx && /^    '"$field"':/ {
      sub(/^    '"$field"':\s*/, ""); gsub(/"/, ""); print; exit
    }
    found && /^[^ ]/ && count>idx { exit }
  ' "$file" | xargs
}

yaml_block_list() {
  local block_key="$1" index="$2" field="$3" file="$4"
  awk -v idx="$index" -v fld="$field" '
    /^'"$block_key"':/ { found=1; next }
    found && /^  - / { count++ }
    found && count==idx && /^      - / && in_list { sub(/^      - /, ""); gsub(/"/, ""); print }
    found && count==idx && /^    '"$field"':/ { in_list=1; next }
    found && count==idx && in_list && !/^      - / { in_list=0 }
    found && /^[^ ]/ && count>idx { exit }
  ' "$file" | xargs -L1
}

# ── Interactive Menu ─────────────────────────────────────────────────────────

interactive_setup() {
  banner

  # ── Check prerequisites ──
  if ! $RC auth current &>/dev/null 2>&1; then
    err "Not authenticated. Run: rc auth login --api-key YOUR_KEY"
    exit 1
  fi
  ok "Authenticated"

  local project_id
  project_id=$($RC doctor 2>&1 | grep -oE 'proj[a-zA-Z0-9_]+' | head -1 || true)
  if [[ -z "$project_id" ]]; then
    err "No project configured. Run: rc init --project YOUR_PROJECT_ID"
    exit 1
  fi
  ok "Project: $project_id"

  # ── Step 1: Platform ──
  select_platform

  # ── Step 2: App identifiers ──
  divider
  echo -e "\n${BOLD}App Details${NC}\n"

  local app_name ios_bundle="" android_package=""
  prompt app_name "App name" "My App"

  if [[ "$PLATFORM" == "ios" || "$PLATFORM" == "both" ]]; then
    prompt ios_bundle "iOS Bundle ID" "com.example.myapp"
  fi
  if [[ "$PLATFORM" == "android" || "$PLATFORM" == "both" ]]; then
    prompt android_package "Android Package Name" "com.example.myapp"
  fi

  # Derive base product ID from bundle/package
  local base_id="${ios_bundle:-$android_package}"

  # ── Step 3: Template or Custom ──
  divider
  local template_choice
  template_choice=$(select_template)

  local -a pkg_keys=() pkg_names=() pkg_ios_ids=() pkg_android_ids=()
  local -a pkg_types=() pkg_entitlements=() pkg_prices=()
  local entitlements_input offering_key offering_name

  if [[ "$template_choice" =~ ^[1-5]$ ]]; then
    apply_template "$template_choice" "$base_id"

    entitlements_input="$TEMPLATE_ENTITLEMENTS"
    offering_key="$TEMPLATE_OFFERING_KEY"
    offering_name="$TEMPLATE_OFFERING_NAME"
    pkg_keys=("${TEMPLATE_PKG_KEYS[@]}")
    pkg_names=("${TEMPLATE_PKG_NAMES[@]}")
    pkg_types=("${TEMPLATE_PKG_TYPES[@]}")
    pkg_prices=("${TEMPLATE_PKG_PRICES[@]}")
    pkg_ios_ids=("${TEMPLATE_PKG_IOS_IDS[@]}")
    pkg_android_ids=("${TEMPLATE_PKG_ANDROID_IDS[@]}")
    pkg_entitlements=("${TEMPLATE_PKG_ENTS[@]}")

    echo ""
    ok "Template loaded. You can customize everything below."
    echo ""

    # Show what the template generated and let them edit
    divider
    echo -e "\n${BOLD}Review & Customize${NC}\n"

    echo -e "  ${BOLD}Entitlements:${NC} $entitlements_input"
    if confirm "Change entitlements?"; then
      prompt entitlements_input "Entitlements (comma-separated)" "$entitlements_input"
    fi

    echo ""
    echo -e "  ${BOLD}Offering:${NC} $offering_key ($offering_name)"
    if confirm "Change offering?"; then
      prompt offering_key "Offering lookup key" "$offering_key"
      prompt offering_name "Offering display name" "$offering_name"
    fi

    echo ""
    echo -e "  ${BOLD}Packages:${NC}"
    for i in "${!pkg_keys[@]}"; do
      local price_display="${pkg_prices[$i]:-}"
      echo -e "    ${CYAN}$((i+1)).${NC} ${pkg_names[$i]}  ${DIM}(${pkg_types[$i]})${NC}  ${GREEN}${price_display}${NC}"
      if [[ "$PLATFORM" == "ios" || "$PLATFORM" == "both" ]]; then
        echo -e "       iOS: ${pkg_ios_ids[$i]}"
      fi
      if [[ "$PLATFORM" == "android" || "$PLATFORM" == "both" ]]; then
        echo -e "       Android: ${pkg_android_ids[$i]}"
      fi
    done
    echo ""

    if confirm "Edit any package?"; then
      edit_packages pkg_keys pkg_names pkg_ios_ids pkg_android_ids pkg_types pkg_prices pkg_entitlements "$ios_bundle" "$android_package"
    fi

    if confirm "Add more packages?"; then
      add_packages pkg_keys pkg_names pkg_ios_ids pkg_android_ids pkg_types pkg_prices pkg_entitlements "$ios_bundle" "$android_package" "$entitlements_input"
    fi

  else
    # ── Full Custom Mode ──
    divider
    echo -e "\n${BOLD}Entitlements${NC}\n"
    prompt entitlements_input "Entitlements (comma-separated lookup keys)" "premium"

    divider
    echo -e "\n${BOLD}Offering${NC}\n"
    prompt offering_key "Offering lookup key" "default"
    prompt offering_name "Offering display name" "Default Offering"

    divider
    echo -e "\n${BOLD}Packages${NC}"
    echo -e "${DIM}Define each subscription/purchase tier.${NC}\n"

    IFS=',' read -ra ent_array <<< "$entitlements_input"
    add_packages pkg_keys pkg_names pkg_ios_ids pkg_android_ids pkg_types pkg_prices pkg_entitlements "$ios_bundle" "$android_package" "$entitlements_input"
  fi

  # ── Summary ──
  show_summary \
    "$app_name" "$ios_bundle" "$android_package" \
    "$entitlements_input" "$offering_key" "$offering_name" \
    pkg_keys pkg_names pkg_ios_ids pkg_android_ids pkg_types pkg_prices pkg_entitlements

  echo ""
  if ! confirm "Create everything?"; then
    echo ""
    if confirm "Save as config file for later?"; then
      local save_path
      prompt save_path "Save config to" "setup-offerings.yaml"
      save_config "$save_path" \
        "$app_name" "$ios_bundle" "$android_package" \
        "$entitlements_input" "$offering_key" "$offering_name" \
        pkg_keys pkg_names pkg_ios_ids pkg_android_ids pkg_types pkg_prices pkg_entitlements
      ok "Saved to $save_path"
      info "Run later with: ./scripts/setup-offerings.sh --config $save_path"
    fi
    exit 0
  fi

  echo ""
  execute_setup \
    "$app_name" "$ios_bundle" "$android_package" \
    "$entitlements_input" \
    "$offering_key" "$offering_name" \
    "$(IFS='|'; echo "${pkg_keys[*]}")" \
    "$(IFS='|'; echo "${pkg_names[*]}")" \
    "$(IFS='|'; echo "${pkg_ios_ids[*]}")" \
    "$(IFS='|'; echo "${pkg_android_ids[*]}")" \
    "$(IFS='|'; echo "${pkg_types[*]}")" \
    "$(IFS='|'; echo "${pkg_entitlements[*]}")" \
    "$(IFS='|'; echo "${pkg_prices[*]}")"
}

# ── Package Editing ──────────────────────────────────────────────────────────

edit_packages() {
  local -n _keys=$1 _names=$2 _ios=$3 _android=$4 _types=$5 _prices=$6 _ents=$7
  local ios_bundle="$8" android_package="$9"

  while true; do
    local idx
    prompt idx "Which package number to edit? (1-${#_keys[@]}, or 'done')" "done"
    [[ "$idx" == "done" ]] && break

    local i=$((idx - 1))
    if [[ $i -lt 0 || $i -ge ${#_keys[@]} ]]; then
      warn "Invalid package number."
      continue
    fi

    echo -e "\n  ${BOLD}Editing: ${_names[$i]}${NC}\n"
    prompt "_keys[$i]" "  Lookup key" "${_keys[$i]}"
    prompt "_names[$i]" "  Display name" "${_names[$i]}"
    prompt "_types[$i]" "  Type (subscription/one_time)" "${_types[$i]}"
    prompt "_prices[$i]" "  Price (display reference)" "${_prices[$i]}"
    if [[ "$PLATFORM" == "ios" || "$PLATFORM" == "both" ]]; then
      prompt "_ios[$i]" "  iOS product ID" "${_ios[$i]}"
    fi
    if [[ "$PLATFORM" == "android" || "$PLATFORM" == "both" ]]; then
      prompt "_android[$i]" "  Android product ID" "${_android[$i]}"
    fi
    prompt "_ents[$i]" "  Entitlements (comma-separated)" "${_ents[$i]}"

    ok "Updated ${_names[$i]}"
    echo ""
  done
}

add_packages() {
  local -n _keys=$1 _names=$2 _ios=$3 _android=$4 _types=$5 _prices=$6 _ents=$7
  local ios_bundle="$8" android_package="$9" default_ent="${10:-premium}"

  IFS=',' read -ra ent_arr <<< "$default_ent"
  local first_ent="${ent_arr[0]}"

  local add_more=true
  while $add_more; do
    local pkg_key pkg_name pkg_ios="" pkg_android="" pkg_type pkg_price pkg_ent
    echo ""
    prompt pkg_key "Package lookup key (e.g. monthly, annual, lifetime)"
    prompt pkg_name "Display name" "${pkg_key^}"
    prompt pkg_type "Type (subscription / one_time)" "subscription"
    prompt pkg_price "Price (e.g. \$9.99/month, \$49.99/year, \$4.99)" ""

    if [[ "$PLATFORM" == "ios" || "$PLATFORM" == "both" ]]; then
      prompt pkg_ios "iOS product identifier" "${ios_bundle}.${pkg_key}"
    fi
    if [[ "$PLATFORM" == "android" || "$PLATFORM" == "both" ]]; then
      prompt pkg_android "Android product identifier" "${android_package:-$ios_bundle}.${pkg_key}"
    fi

    prompt pkg_ent "Attach to entitlements (comma-separated)" "$first_ent"

    _keys+=("$pkg_key")
    _names+=("$pkg_name")
    _ios+=("$pkg_ios")
    _android+=("$pkg_android")
    _types+=("$pkg_type")
    _prices+=("$pkg_price")
    _ents+=("$pkg_ent")

    ok "Added: $pkg_name ($pkg_price)"
    echo ""
    if ! confirm "Add another package?"; then
      add_more=false
    fi
  done
}

# ── Summary Display ──────────────────────────────────────────────────────────

show_summary() {
  local app_name="$1" ios_bundle="$2" android_package="$3"
  local entitlements="$4" offering_key="$5" offering_name="$6"
  local -n s_keys=$7 s_names=$8 s_ios=$9 s_android=${10} s_types=${11} s_prices=${12} s_ents=${13}

  divider
  echo -e "\n${BOLD}${MAGENTA}Setup Summary${NC}\n"

  echo -e "  ${BOLD}Platform:${NC}      $PLATFORM"
  echo -e "  ${BOLD}App:${NC}           $app_name"
  [[ -n "$ios_bundle" ]]       && echo -e "  ${BOLD}iOS:${NC}           $ios_bundle"
  [[ -n "$android_package" ]]  && echo -e "  ${BOLD}Android:${NC}       $android_package"
  echo -e "  ${BOLD}Entitlements:${NC}  $entitlements"
  echo -e "  ${BOLD}Offering:${NC}      $offering_key ($offering_name)"
  echo ""

  # Package table
  local col1=20 col2=14 col3=14 col4=42
  printf "  ${BOLD}%-${col1}s %-${col2}s %-${col3}s %s${NC}\n" "Package" "Type" "Price" "Product IDs"
  printf "  ${DIM}%-${col1}s %-${col2}s %-${col3}s %s${NC}\n" "───────────────────" "─────────────" "─────────────" "─────────────────────────────────────────"

  for i in "${!s_keys[@]}"; do
    local ids=""
    [[ -n "${s_ios[$i]}" ]] && ids+="iOS:${s_ios[$i]}"
    [[ -n "${s_android[$i]}" ]] && { [[ -n "$ids" ]] && ids+=" | "; ids+="Android:${s_android[$i]}"; }

    printf "  %-${col1}s %-${col2}s ${GREEN}%-${col3}s${NC} ${DIM}%s${NC}\n" \
      "${s_names[$i]}" "${s_types[$i]}" "${s_prices[$i]:-—}" "$ids"
  done

  echo ""
  local total_products=0
  for i in "${!s_keys[@]}"; do
    [[ -n "${s_ios[$i]}" ]] && ((total_products++))
    [[ -n "${s_android[$i]}" ]] && ((total_products++))
  done

  local apps=0
  [[ -n "$ios_bundle" ]] && ((apps++))
  [[ -n "$android_package" ]] && ((apps++))

  IFS=',' read -ra ent_list <<< "$entitlements"
  echo -e "  ${DIM}Will create: ${apps} app(s), ${total_products} product(s), ${#ent_list[@]} entitlement(s), 1 offering, ${#s_keys[@]} package(s)${NC}"
  divider
}

# ── Save Config ──────────────────────────────────────────────────────────────

save_config() {
  local file="$1" app_name="$2" ios_bundle="$3" android_package="$4"
  local entitlements="$5" offering_key="$6" offering_name="$7"
  local -n c_keys=$8 c_names=$9 c_ios=${10} c_android=${11} c_types=${12} c_prices=${13} c_ents=${14}

  {
    echo "# Generated by setup-offerings.sh"
    echo "# Run: ./scripts/setup-offerings.sh --config $file"
    echo ""
    echo "platform: \"$PLATFORM\""
    echo "app_name: \"$app_name\""
    [[ -n "$ios_bundle" ]]      && echo "ios_bundle_id: \"$ios_bundle\""
    [[ -n "$android_package" ]] && echo "android_package_name: \"$android_package\""
    echo ""
    echo "entitlements:"
    IFS=',' read -ra ent_list <<< "$entitlements"
    for e in "${ent_list[@]}"; do
      e=$(echo "$e" | xargs)
      echo "  - lookup_key: $e"
      echo "    display_name: \"${e^}\""
    done
    echo ""
    echo "offering:"
    echo "  lookup_key: $offering_key"
    echo "  display_name: \"$offering_name\""
    echo "  set_current: true"
    echo ""
    echo "packages:"
    for i in "${!c_keys[@]}"; do
      echo "  - lookup_key: ${c_keys[$i]}"
      echo "    display_name: \"${c_names[$i]}\""
      [[ -n "${c_ios[$i]}" ]]     && echo "    ios_product_id: \"${c_ios[$i]}\""
      [[ -n "${c_android[$i]}" ]] && echo "    android_product_id: \"${c_android[$i]}\""
      echo "    type: ${c_types[$i]}"
      [[ -n "${c_prices[$i]}" ]]  && echo "    price: \"${c_prices[$i]}\""
      echo "    entitlements:"
      IFS=',' read -ra pkg_ent_list <<< "${c_ents[$i]}"
      for pe in "${pkg_ent_list[@]}"; do
        echo "      - $(echo "$pe" | xargs)"
      done
      echo ""
    done
  } > "$file"
}

# ── Config File Mode ─────────────────────────────────────────────────────────

config_setup() {
  local file="$1"

  if [[ ! -f "$file" ]]; then
    err "Config file not found: $file"
    info "Generate a sample: ./scripts/setup-offerings.sh --sample > setup-offerings.yaml"
    exit 1
  fi

  PLATFORM="$(yaml_val "platform" "$file" 2>/dev/null || echo "both")"
  local app_name ios_bundle android_package offering_key offering_name
  app_name="$(yaml_val "app_name" "$file")"
  ios_bundle="$(yaml_val "ios_bundle_id" "$file" 2>/dev/null || echo "")"
  android_package="$(yaml_val "android_package_name" "$file" 2>/dev/null || echo "")"

  offering_key=$(awk '/^offering:/{found=1} found && /lookup_key:/{sub(/.*lookup_key:\s*/, ""); gsub(/"/, ""); print; exit}' "$file" | xargs)
  offering_name=$(awk '/^offering:/{found=1} found && /display_name:/{sub(/.*display_name:\s*/, ""); gsub(/"/, ""); print; exit}' "$file" | xargs)

  local entitlements_csv
  entitlements_csv=$(awk '/^entitlements:/{found=1; next} found && /^  - /{next} found && /lookup_key:/{sub(/.*lookup_key:\s*/, ""); gsub(/"/, ""); printf "%s,", $0} found && /^[^ ]/{exit}' "$file" | sed 's/,$//')

  local pkg_count
  pkg_count=$(yaml_block_count "packages" "$file")

  local pkg_keys="" pkg_names="" pkg_ios="" pkg_android="" pkg_types="" pkg_ents="" pkg_prices=""
  for ((i=1; i<=pkg_count; i++)); do
    local k n ios and typ ent price
    k=$(yaml_block_field "packages" "$i" "lookup_key" "$file")
    n=$(yaml_block_field "packages" "$i" "display_name" "$file")
    ios=$(yaml_block_field "packages" "$i" "ios_product_id" "$file" 2>/dev/null || echo "")
    and=$(yaml_block_field "packages" "$i" "android_product_id" "$file" 2>/dev/null || echo "")
    typ=$(yaml_block_field "packages" "$i" "type" "$file")
    price=$(yaml_block_field "packages" "$i" "price" "$file" 2>/dev/null || echo "")
    ent=$(yaml_block_list "packages" "$i" "entitlements" "$file" | paste -sd, -)

    [[ -n "$pkg_keys" ]]   && pkg_keys+="|"
    [[ -n "$pkg_names" ]]  && pkg_names+="|"
    [[ -n "$pkg_ios" ]]    && pkg_ios+="|"
    [[ -n "$pkg_android" ]] && pkg_android+="|"
    [[ -n "$pkg_types" ]]  && pkg_types+="|"
    [[ -n "$pkg_ents" ]]   && pkg_ents+="|"
    [[ -n "$pkg_prices" ]] && pkg_prices+="|"

    pkg_keys+="$k"; pkg_names+="$n"; pkg_ios+="$ios"; pkg_android+="$and"
    pkg_types+="$typ"; pkg_ents+="$ent"; pkg_prices+="$price"
  done

  echo -e "\n${BOLD}${CYAN}RevenueCat Setup from Config${NC}"
  echo -e "${DIM}$file${NC}\n"
  echo "  Platform:      $PLATFORM"
  echo "  App:           $app_name"
  [[ -n "$ios_bundle" ]]      && echo "  iOS:           $ios_bundle"
  [[ -n "$android_package" ]] && echo "  Android:       $android_package"
  echo "  Entitlements:  $entitlements_csv"
  echo "  Offering:      $offering_key ($offering_name)"
  echo "  Packages:      $pkg_count"
  echo ""

  if ! $DRY_RUN && ! confirm "Proceed?"; then
    exit 0
  fi

  echo ""
  execute_setup \
    "$app_name" "$ios_bundle" "$android_package" \
    "$entitlements_csv" \
    "$offering_key" "$offering_name" \
    "$pkg_keys" "$pkg_names" "$pkg_ios" "$pkg_android" "$pkg_types" "$pkg_ents" "$pkg_prices"
}

# ── Execution Engine ─────────────────────────────────────────────────────────

execute_setup() {
  local app_name="$1" ios_bundle="$2" android_package="$3"
  local entitlements_csv="$4"
  local offering_key="$5" offering_name="$6"
  local pkg_keys_str="$7" pkg_names_str="$8"
  local pkg_ios_str="$9" pkg_android_str="${10}"
  local pkg_types_str="${11}" pkg_ents_str="${12}"
  local pkg_prices_str="${13:-}"

  IFS=',' read -ra ent_keys <<< "$entitlements_csv"
  IFS='|' read -ra pkg_keys <<< "$pkg_keys_str"
  IFS='|' read -ra pkg_names <<< "$pkg_names_str"
  IFS='|' read -ra pkg_ios <<< "$pkg_ios_str"
  IFS='|' read -ra pkg_android <<< "$pkg_android_str"
  IFS='|' read -ra pkg_types <<< "$pkg_types_str"
  IFS='|' read -ra pkg_ents <<< "$pkg_ents_str"
  IFS='|' read -ra pkg_prices <<< "$pkg_prices_str"

  TOTAL_STEPS=6

  local do_ios=false do_android=false
  [[ "$PLATFORM" == "ios" || "$PLATFORM" == "both" ]] && [[ -n "$ios_bundle" ]] && do_ios=true
  [[ "$PLATFORM" == "android" || "$PLATFORM" == "both" ]] && [[ -n "$android_package" ]] && do_android=true

  # ── Step 1: Create Apps ──
  step 1 "Creating app(s)"

  if $do_ios; then
    local ios_output
    ios_output=$(run_rc apps create --name "${app_name} iOS" --type app_store --bundle-id "$ios_bundle" -o json 2>&1) || true
    IOS_APP_ID=$(echo "$ios_output" | extract_id)
    if [[ -n "$IOS_APP_ID" ]]; then
      ok "iOS app: $IOS_APP_ID"
    elif $DRY_RUN; then
      IOS_APP_ID="app_ios_dry_run"
      ok "(dry-run) iOS app"
    else
      warn "iOS app creation: $ios_output"
      prompt IOS_APP_ID "Enter existing iOS app ID (or press enter to skip)" ""
    fi
  fi

  if $do_android; then
    local android_output
    android_output=$(run_rc apps create --name "${app_name} Android" --type play_store --package-name "$android_package" -o json 2>&1) || true
    ANDROID_APP_ID=$(echo "$android_output" | extract_id)
    if [[ -n "$ANDROID_APP_ID" ]]; then
      ok "Android app: $ANDROID_APP_ID"
    elif $DRY_RUN; then
      ANDROID_APP_ID="app_android_dry_run"
      ok "(dry-run) Android app"
    else
      warn "Android app creation: $android_output"
      prompt ANDROID_APP_ID "Enter existing Android app ID (or press enter to skip)" ""
    fi
  fi

  # ── Step 2: Create Products ──
  step 2 "Creating products"

  declare -A ios_product_ids=()
  declare -A android_product_ids=()

  for i in "${!pkg_keys[@]}"; do
    local key="${pkg_keys[$i]}"
    local type="${pkg_types[$i]}"
    local price_label="${pkg_prices[$i]:-}"

    # iOS product
    if $do_ios && [[ -n "${pkg_ios[$i]}" && -n "$IOS_APP_ID" ]]; then
      local ios_prod_out
      ios_prod_out=$(run_rc products create --store-identifier "${pkg_ios[$i]}" --type "$type" --app-id "$IOS_APP_ID" -o json 2>&1) || true
      local ios_pid=$(echo "$ios_prod_out" | extract_id)
      if [[ -n "$ios_pid" ]]; then
        ios_product_ids["$key"]="$ios_pid"
        ok "  iOS [$key]: $ios_pid  ${GREEN}${price_label}${NC}"
      elif $DRY_RUN; then
        ios_product_ids["$key"]="prod_ios_${key}_dry"
        ok "  (dry-run) iOS [$key]  ${GREEN}${price_label}${NC}"
      else
        warn "  iOS [$key] failed: $ios_prod_out"
      fi
    fi

    # Android product
    if $do_android && [[ -n "${pkg_android[$i]}" && -n "$ANDROID_APP_ID" ]]; then
      local and_prod_out
      and_prod_out=$(run_rc products create --store-identifier "${pkg_android[$i]}" --type "$type" --app-id "$ANDROID_APP_ID" -o json 2>&1) || true
      local and_pid=$(echo "$and_prod_out" | extract_id)
      if [[ -n "$and_pid" ]]; then
        android_product_ids["$key"]="$and_pid"
        ok "  Android [$key]: $and_pid  ${GREEN}${price_label}${NC}"
      elif $DRY_RUN; then
        android_product_ids["$key"]="prod_and_${key}_dry"
        ok "  (dry-run) Android [$key]  ${GREEN}${price_label}${NC}"
      else
        warn "  Android [$key] failed: $and_prod_out"
      fi
    fi
  done

  # ── Step 3: Create Entitlements ──
  step 3 "Creating entitlements"

  declare -A entitlement_ids=()

  for ent_key in "${ent_keys[@]}"; do
    ent_key=$(echo "$ent_key" | xargs)
    local display="${ent_key^}"
    local ent_out
    ent_out=$(run_rc entitlements create --lookup-key "$ent_key" --display-name "$display" -o json 2>&1) || true
    local ent_id=$(echo "$ent_out" | extract_id)
    if [[ -n "$ent_id" ]]; then
      entitlement_ids["$ent_key"]="$ent_id"
      ok "  Entitlement [$ent_key]: $ent_id"
    elif $DRY_RUN; then
      entitlement_ids["$ent_key"]="entl_${ent_key}_dry"
      ok "  (dry-run) Entitlement [$ent_key]"
    else
      warn "  Entitlement [$ent_key]: $ent_out"
    fi
  done

  # ── Step 4: Attach Products to Entitlements ──
  step 4 "Attaching products to entitlements"

  for i in "${!pkg_keys[@]}"; do
    local key="${pkg_keys[$i]}"
    IFS=',' read -ra pkg_ent_list <<< "${pkg_ents[$i]}"

    for ent_key in "${pkg_ent_list[@]}"; do
      ent_key=$(echo "$ent_key" | xargs)
      local ent_id="${entitlement_ids[$ent_key]:-}"
      [[ -z "$ent_id" ]] && { warn "  Skipping: entitlement '$ent_key' not found"; continue; }

      local attach_ids=""
      [[ -n "${ios_product_ids[$key]:-}" ]] && attach_ids="${ios_product_ids[$key]}"
      [[ -n "${android_product_ids[$key]:-}" ]] && {
        [[ -n "$attach_ids" ]] && attach_ids+=","
        attach_ids+="${android_product_ids[$key]}"
      }

      if [[ -n "$attach_ids" ]]; then
        run_rc entitlements attach-products --entitlement-id "$ent_id" --product-ids "$attach_ids" >/dev/null 2>&1 || true
        ok "  [$key] -> [$ent_key]"
      fi
    done
  done

  # ── Step 5: Create Offering + Packages ──
  step 5 "Creating offering and packages"

  local off_out
  off_out=$(run_rc offerings create --lookup-key "$offering_key" --display-name "$offering_name" -o json 2>&1) || true
  OFFERING_ID=$(echo "$off_out" | extract_id)
  if [[ -n "$OFFERING_ID" ]]; then
    ok "Offering: $OFFERING_ID ($offering_key)"
  elif $DRY_RUN; then
    OFFERING_ID="ofrngs_dry_run"
    ok "(dry-run) Offering"
  else
    warn "Offering: $off_out"
    prompt OFFERING_ID "Enter existing offering ID" ""
  fi

  if [[ -n "$OFFERING_ID" ]]; then
    run_rc offerings update --offering-id "$OFFERING_ID" --is-current >/dev/null 2>&1 || true
    ok "Set as current offering"
  fi

  for i in "${!pkg_keys[@]}"; do
    local key="${pkg_keys[$i]}" name="${pkg_names[$i]}"

    local pkg_out
    pkg_out=$(run_rc packages create --offering-id "$OFFERING_ID" --lookup-key "$key" --display-name "$name" -o json 2>&1) || true
    local pkg_id=$(echo "$pkg_out" | extract_id)

    if [[ -n "$pkg_id" ]]; then
      ok "  Package [$key]: $pkg_id"
      PACKAGE_IDS+=("$pkg_id")
    elif $DRY_RUN; then
      pkg_id="pkg_${key}_dry"
      ok "  (dry-run) Package [$key]"
      PACKAGE_IDS+=("$pkg_id")
    else
      warn "  Package [$key]: $pkg_out"
      continue
    fi

    local attach_ids=""
    [[ -n "${ios_product_ids[$key]:-}" ]] && attach_ids="${ios_product_ids[$key]}"
    [[ -n "${android_product_ids[$key]:-}" ]] && {
      [[ -n "$attach_ids" ]] && attach_ids+=","
      attach_ids+="${android_product_ids[$key]}"
    }

    if [[ -n "$attach_ids" && -n "$pkg_id" ]]; then
      run_rc packages attach-products --package-id "$pkg_id" --product-ids "$attach_ids" >/dev/null 2>&1 || true
      ok "  Attached products to [$key]"
    fi
  done

  # ── Step 6: Done ──
  step 6 "Complete"
  echo ""
  divider
  echo -e "\n${BOLD}${GREEN}Setup Complete!${NC}\n"

  local total_apps=0 total_products=0
  [[ -n "$IOS_APP_ID" ]] && ((total_apps++))
  [[ -n "$ANDROID_APP_ID" ]] && ((total_apps++))
  for key in "${pkg_keys[@]}"; do
    [[ -n "${ios_product_ids[$key]:-}" ]] && ((total_products++))
    [[ -n "${android_product_ids[$key]:-}" ]] && ((total_products++))
  done

  echo -e "  ${BOLD}Created:${NC}"
  echo "    Apps:          $total_apps"
  [[ -n "$IOS_APP_ID" ]]      && echo "      iOS:         $IOS_APP_ID"
  [[ -n "$ANDROID_APP_ID" ]]  && echo "      Android:     $ANDROID_APP_ID"
  echo "    Products:      $total_products"
  echo "    Entitlements:  ${#entitlement_ids[@]}"
  echo "    Offering:      $OFFERING_ID"
  echo "    Packages:      ${#PACKAGE_IDS[@]}"

  echo ""
  echo -e "  ${BOLD}Price Reference:${NC}"
  for i in "${!pkg_keys[@]}"; do
    local p="${pkg_prices[$i]:-—}"
    printf "    %-20s %s\n" "${pkg_names[$i]}" "$p"
  done

  echo ""
  echo -e "  ${DIM}Verify:${NC}  rc status"
  echo -e "  ${DIM}View:${NC}    rc offerings get --offering-id $OFFERING_ID"
  echo ""
  divider
}

# ── Main ─────────────────────────────────────────────────────────────────────

usage() {
  cat <<EOF
Usage: setup-offerings.sh [options]

One-click RevenueCat setup — creates apps, products, entitlements, offerings,
and packages for iOS, Android, or both platforms.

Options:
  --config FILE    Use a YAML config file
  --dry-run        Preview without creating anything
  --sample         Print a sample config file
  --verbose        Show all rc commands
  --help           Show this help

Examples:
  ./scripts/setup-offerings.sh                               # Interactive menu
  ./scripts/setup-offerings.sh --config plan.yaml            # From config
  ./scripts/setup-offerings.sh --dry-run --config plan.yaml  # Preview
  ./scripts/setup-offerings.sh --sample > plan.yaml          # Generate template
EOF
}

main() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --config)   CONFIG_FILE="$2"; shift 2 ;;
      --dry-run)  DRY_RUN=true; shift ;;
      --sample)   generate_sample_config; exit 0 ;;
      --verbose)  VERBOSE=true; shift ;;
      --help|-h)  usage; exit 0 ;;
      *)          err "Unknown option: $1"; usage; exit 1 ;;
    esac
  done

  if ! command -v "$RC" &>/dev/null; then
    err "rc (revenuecat-cli) not found. Install:"
    echo "  brew tap AndroidPoet/tap && brew install revenuecat-cli"
    exit 1
  fi

  if $DRY_RUN; then
    warn "DRY RUN — no resources will be created"
    echo ""
  fi

  if [[ -n "$CONFIG_FILE" ]]; then
    config_setup "$CONFIG_FILE"
  else
    interactive_setup
  fi
}

main "$@"

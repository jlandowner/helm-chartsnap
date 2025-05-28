#!/bin/sh
set -e

[ "$HELM_DEBUG" != "false" ] && set -x && (printenv | grep HELM)


# Function to print error message and exit
error_exit() {
    echo "$1" >&2
    exit 1
}

# Function to validate command availability
validate_command() {
    command -v "$1" >/dev/null 2>&1 || error_exit "Required command '$1' not found. Please install it."
}

# Function to detect the specified version from plugin.yaml
get_plugin_version() {
    cat $HELM_PLUGIN_DIR/plugin.yaml | grep version: | cut -d " " -f 2
}

# Function to download and install the plugin
install_plugin() {
    local plugin_version="$1"
    local plugin_url="$2"
    local plugin_filename="$3"
    local plugin_directory="$4"

    # Download the plugin archive
    if validate_command "curl"; then
        curl --fail -sSL "${plugin_url}" -o "${plugin_filename}"
    elif validate_command "wget"; then
        wget -q "${plugin_url}" -O "${plugin_filename}"
    else
        error_exit "Both 'curl' and 'wget' commands not found. Please install either one."
    fi

    # Extract and install the plugin
    tar xzf "${plugin_filename}" -C "${plugin_directory}"
    mv "${plugin_directory}/${name}" "bin/${name}" || mv "${plugin_directory}/${name}.exe" "bin/${name}"
}

# Main script
name="chartsnap"
repo_name="helm-chartsnap"
repo_owner="jlandowner"
repo="https://github.com/${repo_owner}/${repo_name}"
HELM_PUSH_PLUGIN_NO_INSTALL_HOOK="${HELM_PUSH_PLUGIN_NO_INSTALL_HOOK:-}"

# Check if in development mode
if [ -n "$HELM_PUSH_PLUGIN_NO_INSTALL_HOOK" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

# If update flag is provided, git checkout the latest tag
if [ "$1" = "-u" ]; then
    echo "Updating ${name} plugin..."
    git fetch --tags || error_exit "Failed to fetch tags from remote repository"
    
    latest_tag=$(git tag --sort=-creatordate | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1)
    current_tag=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ "$current_tag" = "$latest_tag" ]; then
        echo "${name} is already up to date (${latest_tag})."
        exit 0
    fi
    git checkout "$latest_tag" || error_exit "Failed to checkout: $latest_tag"
fi

# Autodetect the plugin version
version=$(get_plugin_version)
echo "Downloading and installing ${name} v${version} ..."

# Convert architecture of the target system to a compatible GOARCH value
case $(uname -m) in
    x86_64)
        arch="amd64"
        ;;
    aarch64 | arm64)
        arch="arm64"
        ;;
    *)
        error_exit "Failed to detect target architecture"
        ;;
esac

# Construct the plugin download URL
if [ "$(uname)" = "Darwin" ]; then
    url="${repo}/releases/download/v${version}/${name}_v${version}_darwin_${arch}.tar.gz"
elif [ "$(uname)" = "Linux" ] ; then
    url="${repo}/releases/download/v${version}/${name}_v${version}_linux_${arch}.tar.gz"
else
    url="${repo}/releases/download/v${version}/${name}_v${version}_windows_${arch}.tar.gz"
fi

echo "$url"

mkdir -p "bin"
mkdir -p "releases/${version}"

install_plugin "$version" "$url" "releases/${version}.tar.gz" "releases/${version}"

echo
echo "${name} is installed. To start it, run the following command:"
echo "  helm chartsnap"
echo

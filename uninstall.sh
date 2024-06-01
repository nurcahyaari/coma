#!/bin/bash

# Variables
TOOL_NAME="coma"
TOOL_BINARY="/usr/local/bin/$TOOL_NAME"
CONFIG_DIR="/usr/local/opt/$TOOL_NAME"
DATA_DIR="/var/lib/$TOOL_NAME"
LOG_DIR="/var/log/coma"
USER="comauser"
SERVICE_FILE="/etc/systemd/system/$TOOL_NAME.service"
PLIST_PUBLISHER="com.nurcahyaari.$TOOL_NAME"
PLIST_FILE="$PLIST_PUBLISHER.plist"
LAUNCHD_DIR="/Library/LaunchDaemons"

# Linux
remove_systemd_service() {
    echo "Removing systemd service..."
    sudo systemctl stop $TOOL_NAME
    sudo systemctl disable $TOOL_NAME
    sudo rm -f $SERVICE_FILE
    sudo systemctl daemon-reload
}

# macOS
remove_launchd_service() {
    echo "Removing launchd service..."
    sudo launchctl stop $PLIST_PUBLISHER
    sudo launchctl unload $LAUNCHD_DIR/$PLIST_FILE
    sudo rm -f $LAUNCHD_DIR/$PLIST_FILE
}

remove_files() {
    echo "Removing binary and directories..."
    sudo rm -f $TOOL_BINARY
    sudo rm -rf $CONFIG_DIR
    sudo rm -rf $DATA_DIR
    sudo rm -rf $LOG_DIR
}

remove_user() {
    if id -u $USER > /dev/null 2>&1; then
        echo "Removing user $USER..."
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            sudo userdel $USER
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            sudo dscl . -delete /Users/$USER
        fi
    fi
}

main() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        remove_systemd_service
        echo "Uninstallation completed for systemd."
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        remove_launchd_service
        echo "Uninstallation completed for launchd."
    else
        echo "Unsupported OS."
        exit 1
    fi

    remove_files
    remove_user

    echo "$TOOL_NAME has been uninstalled."
}

main

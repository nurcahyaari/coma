#!/bin/bash

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

build_binary() {
    echo "Building the Golang binary..."
    go build -o $TOOL_NAME
    if [ $? -ne 0 ]; then
        echo "Failed to build the binary."
        exit 1
    fi
}

create_user() {
    if ! id -u $USER &> /dev/null; then
        echo "Creating system user $USER..."
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            sudo useradd -r -s /bin/false $USER
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            sudo dscl . -create /Users/$USER
            sudo dscl . -create /Users/$USER UserShell /bin/bash
            sudo dscl . -create /Users/$USER UniqueID "510"
            sudo dscl . -create /Users/$USER PrimaryGroupID 20
        fi
    fi
}

setup_directories() {
    echo "Setting up directories..."
    sudo mkdir -p $CONFIG_DIR
    sudo mkdir -p $DATA_DIR
    sudo mkdir -p $DATA_DIR/database
    sudo mkdir -p $LOG_DIR

    echo "Setting ownership and permissions of the directories..."
    sudo chown -R $USER $CONFIG_DIR
    sudo chmod 750 $CONFIG_DIR
    sudo chown -R $USER $DATA_DIR
    sudo chmod 750 $DATA_DIR
    sudo chown -R $USER $LOG_DIR
    sudo chmod 750 $LOG_DIR
}

copy_binary() {
    echo "Copying the binary to $TOOL_BINARY..."
    sudo mv $TOOL_NAME $TOOL_BINARY
    sudo chmod +x $TOOL_BINARY
    sudo chown -R $USER $TOOL_BINARY
}

# Linux
setup_systemd() {
    echo "Setting up systemd service..."
    cat <<EOL | sudo tee $SERVICE_FILE
[Unit]
Description=Coma HTTP Service
After=network.target

[Service]
ExecStart=$TOOL_BINARY
WorkingDirectory=$DATA_DIR
User=$USER
Group=$USER
Restart=always
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Environment=GO_ENV=production

[Install]
WantedBy=multi-user.target
EOL

    sudo systemctl daemon-reload
    sudo systemctl enable $TOOL_NAME
    sudo systemctl start $TOOL_NAME
}

# macOS
setup_launchd() {
    echo "Setting up launchd service..."
    cat <<EOL | sudo tee $LAUNCHD_DIR/$PLIST_FILE
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$PLIST_PUBLISHER</string>
    <key>ProgramArguments</key>
    <array>
        <string>$TOOL_BINARY</string>
    </array>
    <key>WorkingDirectory</key>
    <string>$DATA_DIR</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/coma/coma_stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/coma/coma_stderr.log</string>
</dict>
</plist>
EOL

    sudo launchctl load $LAUNCHD_DIR/$PLIST_FILE
    sudo launchctl start $PLIST_PUBLISHER
}

main() {
    build_binary
    create_user
    setup_directories
    copy_binary

    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        setup_systemd
        echo "Installation completed. $TOOL_NAME is running as a systemd service."
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        setup_launchd
        echo "Installation completed. $TOOL_NAME is running as a launchd service."
    else
        echo "Unsupported OS."
        exit 1
    fi

    echo "Configuration files should be placed at $CONFIG_DIR."
    echo "Data files should be stored at $DATA_DIR."
}

main

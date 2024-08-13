[Setup]
AppName=Vaulty Remote Clipboard
AppVersion=0.1.0
DefaultDirName={commonpf}\Vaulty Remote Clipboard
DefaultGroupName=Vaulty Remote Clipboard
OutputBaseFilename=VaultyRemoteClipboardSetup
Compression=lzma
SolidCompression=yes

[Files]
Source: "fyne-cross\bin\windows-amd64\Vaulty Remote Clipboard.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Vaulty Remote Clipboard"; Filename: "{app}\Vaulty Remote Clipboard.exe"
Name: "{group}\Uninstall Vaulty Remote Clipboard"; Filename: "{uninstallexe}"

[Run]
Filename: "{app}\Vaulty Remote Clipboard.exe"; Description: "Launch Vaulty Remote Clipboard"; Flags: nowait postinstall skipifsilent

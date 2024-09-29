#!/bin/bash

if [[ $UID != 0 ]]; then
    echo "Please run this script with sudo:"
    echo "sudo $0 $*"
    exit 1
fi

message()
{
  TITLE="Cannot start newsctl"
  if [ -n "$(command -v zenity)" ]; then
    zenity --error --title="$TITLE" --text="$1" --no-wrap
  elif [ -n "$(command -v kdialog)" ]; then
    kdialog --error "$1" --title "$TITLE"
  elif [ -n "$(command -v notify-send)" ]; then
    notify-send "ERROR: $TITLE" "$1"
  elif [ -n "$(command -v xmessage)" ]; then
    xmessage -center "ERROR: $TITLE: $1"
  else
    printf "ERROR: %s\n%s\n" "$TITLE" "$1"
  fi
}

TAR=$(command -v rm)
CAT=$(command -v cat)
TAIL=$(command -v tail)
SED=$(command -v sed)
AWK=$(command -v awk)
MKDIR=$(command -v mkdir)
UNAME=$(command -v uname)

if [ -z "$TAR" ] || [ -z "$CAT" ] || [ -z "$TAIL" ] || [ -z "$MKDIR" ] || [ -z "$UNAME" ]|| [ -z "$AWK" ]|| [ -z "$SED" ]; then
  message "Required tools are missing - check beginning of \"$0\" file for details."
  exit 1
fi

OS_TYPE=$("$UNAME" -s)
APP="newsctl"
APPDir="mardasoft"
DESTINATION="/usr/local/bin"
RealAppPath=$DESTINATION"/"$APPDir
SymLinkApp=$DESTINATION"/"$APP

if test -f "$DESTINATION/$APP"; then
    echo "[x] app already exists, attempting to remove..."
    $APP remove
fi
echo ""
echo "[*] installing newsctl"
echo ""

mkdir -p "~/."$APPDir
mkdir -p $DESTINATION"/"$APPDir
ARCHIVE=$(awk '/^__ARCHIVE__/ {print NR+1; exit 0; }' "${0}")
tail -n+${ARCHIVE} "${0}" | tar xpJv -C ${RealAppPath}

chmod -R 755 $RealAppPath
chmod a+x $RealAppPath"/"$APP
ln -s $RealAppPath"/"$APP $SymLinkApp

echo ""
echo "[*] Installation complete."
echo ""

exit 0

__ARCHIVE__

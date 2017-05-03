#!/bin/bash

DEFAULTSGO="./defaults.go"
HTMLTEMPLATE="./tmpl/gotable.tmpl"
DEFAULTCSS="./gotable.css"

# check for html template file existance
if [ ! -f ${HTMLTEMPLATE} ]; then
    echo "Default HTML Template not found for gotable!"
    exit 0
fi

# get the template content
HTMLTEMPLATE=$(cat ${HTMLTEMPLATE})

# check for default css file existance
if [ ! -f ${DEFAULTCSS} ]; then
    echo "Default CSS file not found for gotable!"
    exit 0
fi

# get the css content
DEFAULTCSS=$(cat ${DEFAULTCSS})

# -----------------------------------------------------------------
# Here document containing the body of the generated script
# -----------------------------------------------------------------
cat  >$DEFAULTSGO <<EOF
package gotable

// DCSS et. al. are the constants used for default values
const (
    DCSS = "${DEFAULTCSS}"
    DTEMPLATE = "${HTMLTEMPLATE}"
)
EOF

# -----------------------------------------------------------------
# if file created successfully then go format it
# -----------------------------------------------------------------
if [ -f "$DEFAULTSGO" ]
then
  # go format it
  gofmt -s -w $DEFAULTSGO
else
  echo "Problem in creating file: \"$DEFAULTSGO\""
fi

exit 0

# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2020, Filippov Alex
#
# This library is free software: you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 3 of the License, or (at your option) any later version.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Library General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public
# License along with this library.  If not, see
# <https://www.gnu.org/licenses/>.

#!/usr/bin/env bash

# hab https://github.com/go-swagger/go-swagger
# docs https://goswagger.io/generate/spec.html

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
SWAGGER=swagger_darwin_amd64.dms.v0.19

__validate() {
    ${SWAGGER} validate ${ROOT}/api/server/docs/swagger.yaml
}

__swagger() {
    cd ${ROOT}/api/server
    ${SWAGGER} generate spec -o ${ROOT}/api/server/docs/swagger.yaml --scan-models
}

main() {

    case "$1" in
        swagger)
        __swagger
        ;;
        *)
        __help
        exit 1
        ;;
    esac

}

__help() {
  cat <<EOF
Usage: generator.sh [options]

OPTIONS:

  swagger - generate an API server

  -h / --help - show this help text and exit 0

EOF
}

main "$@"

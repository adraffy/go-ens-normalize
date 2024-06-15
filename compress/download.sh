#! /usr/bin/env sh

# download spec
curl -o ./data/spec.json https://raw.githubusercontent.com/adraffy/ens-normalize.js/main/derive/output/spec.json
curl -o ./data/nf.json https://raw.githubusercontent.com/adraffy/ens-normalize.js/main/derive/output/nf.json

# download tests
curl -o ./data/tests.json https://raw.githubusercontent.com/adraffy/ens-normalize.js/main/validate/tests.json
curl -o ./data/nf-tests.json https://raw.githubusercontent.com/adraffy/ens-normalize.js/main/derive/output/nf-tests.json

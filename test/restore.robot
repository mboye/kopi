** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py
Resource    common.robot

Test Setup     Begin test
Test Teardown  End test

** Variables **
${small file}           test/resources/restore/small-file.txt
${small file hash}      144062aa1d1186d6ef1c122d645b567a
${large file}           test/resources/restore/large-file.txt
${large file hash 1}    ${small file hash}
${large file hash 2}    074e8e431cc1335d6a44f366adf0eb11
${large file hash}      2f0f639c17a26a374e5063bcd46f5146
${source dir}           test/resources/restore
${store dir}            ${TEMPDIR}/simple_store_data
${restore dir}          ${TEMPDIR}/restored_data
${index}                ${TEMPDIR}/index
${stored index}         ${TEMPDIR}/index.stored
${max block size}       64

** Test Cases **
File of same size as block size
    Create index from "${small file}" and save it to "${index}"
    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have md5 hash   ${restore dir}/${small file}  ${small file hash}

File larger than block size
    Create index from "${large file}" and save it to "${index}"
    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${large file}
    File should have md5 hash   ${restore dir}/${large file}  ${large file hash}

Files reuse existing blocks
    Create index from "${source dir}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  3

    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have md5 hash   ${restore dir}/${small file}  ${small file hash}

    File should exist           ${restore dir}/${large file}
    File should have md5 hash   ${restore dir}/${large file}  ${large file hash}


Dry run
    Create index from "${source dir}" and save it to "${index}"

    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index dry run "${stored index}" from "${store dir}" to "${restore dir}"

    Run keyword and expect error  *
    ...     File should exist     ${restore dir}/${small file}

    Run keyword and expect error  *
    ...     File should exist     ${restore dir}/${large file}

Dry run with missing block
    Create index from "${source dir}" and save it to "${index}"

    Store index "${index}" to "${store dir}" and save output to "${stored index}"

    ${block dir}=       Get substring  ${small file hash}  0  2
    Remove file         ${store dir}/${block dir}/${small file hash}.block

    Run keyword and expect error    *failed to open bloc*
    ...     Restore index dry run "${stored index}" from "${store dir}" to "${restore dir}"

Dry run with block corruption
    Create index from "${source dir}" and save it to "${index}"

    Store index "${index}" to "${store dir}" and save output to "${stored index}"

    ${block dir}=       Get substring  ${small file hash}  0  2
    Create file      ${store dir}/${block dir}/${small file hash}.block   extra-data

    Run keyword and expect error    *corrupt block detected*
    ...     Restore index dry run "${stored index}" from "${store dir}" to "${restore dir}"

** Keywords **
Begin test
    Create directory        ${store dir}
    Create directory        ${restore dir}

End test
    Remove directory  ${store dir}  recursive=True
    Remove directory  ${restore dir}  recursive=True
    Remove file       ${stored index}

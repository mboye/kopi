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
${source dir}           test/resources/store
${restore dir}          ${TEMPDIR}/restored_data
${stored index}         ${TEMPDIR}/index.stored

** Test Cases **
Restore small file
    Create index from "${small file}" and save it to "${index}"
    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have SHA1 hash   ${restore dir}/${small file}  ${small file hash}

Restore small file with encryption
    Create index from "${small file}" and save it to "${index}"
    Store index "${index}" with encryption to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" with encryption from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have SHA1 hash  ${restore dir}/${small file}  ${small file hash}

Restore large file
    Create index from "${large file}" and save it to "${index}"
    Store index "${index}" with encryption to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" with encryption from "${store dir}" to "${restore dir}"

    File should exist            ${restore dir}/${large file}
    File should have SHA1 hash   ${restore dir}/${large file}  ${large file hash}

Restore large file with encryption
    Create index from "${large file}" and save it to "${index}"
    Store index "${index}" with encryption to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" with encryption from "${store dir}" to "${restore dir}"

    File should exist            ${restore dir}/${large file}
    File should have SHA1 hash   ${restore dir}/${large file}  ${large file hash}

Restore multiple files
    Create index from "${source dir}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  3

    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have SHA1 hash   ${restore dir}/${small file}  ${small file hash}

    File should exist           ${restore dir}/${large file}
    File should have SHA1 hash   ${restore dir}/${large file}  ${large file hash}

Restore multiple files with encryption
    Create index from "${source dir}" and save it to "${index}"
    Store index "${index}" with encryption to "${store dir}" and save output to "${stored index}"
    Restore index "${stored index}" with encryption from "${store dir}" to "${restore dir}"

    File should exist           ${restore dir}/${small file}
    File should have SHA1 hash   ${restore dir}/${small file}  ${small file hash}

    File should exist           ${restore dir}/${large file}
    File should have SHA1 hash   ${restore dir}/${large file}  ${large file hash}

Restore dry run
    Create index from "${source dir}" and save it to "${index}"

    Store index "${index}" to "${store dir}" and save output to "${stored index}"
    Restore index dry run "${stored index}" from "${store dir}" to "${restore dir}"

    Run keyword and expect error  *
    ...     File should exist     ${restore dir}/${small file}

    Run keyword and expect error  *
    ...     File should exist     ${restore dir}/${large file}

Restore dry run with missing block
    Create index from "${source dir}" and save it to "${index}"

    Store index "${index}" to "${store dir}" and save output to "${stored index}"

    ${block dir}=       Get substring  ${small file hash}  0  2
    Remove file         ${store dir}/${block dir}/${small file hash}.block

    Run keyword and expect error    *failed to open bloc*
    ...     Restore index dry run "${stored index}" from "${store dir}" to "${restore dir}"

Restore dry run with block corruption
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
    Copy file               test/resources/salt  ${store dir}/salt

End test
    Remove directory  ${store dir}      recursive=True
    Remove directory  ${restore dir}    recursive=True
    Remove file       ${stored index}

** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py

Test Setup     Begin test
Test Teardown  End test

** Variables **
${differ bin}           bin/kopi-diff
${indexer bin}          bin/kopi-index
${store bin}            bin/kopi-store
${restore bin}          bin/kopi-restore
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

** Keywords **
Create index from "${path}" and save it to "${output path}"
    ${rc}=  Run and return RC  ${indexer bin} --init=true ${path} > ${output path} 2>/dev/null
    Should be equal as integers  ${rc}  0

Store index "${index}" to "${store dir}" and save output to "${output path}"
    ${rc}  ${stdout}=  Run and return RC and output  ${store bin} --maxBlockSize ${max block size} ${store dir} < ${index} > ${output path} 2>/dev/null
    Should be equal as integers  ${rc}  0
    Log many  ${stdout}

    ${lines}=   Split to lines  ${stdout}
    [Return]   ${lines}

Restore index "${index}" from "${store dir}" to "${restore dir}"
    ${rc}  ${stdout}=  Run and return RC and output  ${restore bin} ${store dir} ${restore dir} < ${index} 2>&1
    Log many  ${stdout}
    Should be equal as integers  ${rc}  0

Begin test
    Create directory        ${store dir}
    Create directory        ${restore dir}

End test
    Remove directory  ${store dir}  recursive=True
    Remove directory  ${restore dir}  recursive=True
    Remove file       ${stored index}

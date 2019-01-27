** Variables **
${differ bin}           bin/kopi-diff
${indexer bin}          bin/kopi-index
${store bin}            bin/kopi-store
${restore bin}          bin/kopi-restore

${store dir}            ${TEMPDIR}/simple_store_data
${index}                ${TEMPDIR}/index
${max block size}       64

${small file}           test/resources/store/small-file.txt
${large file}           test/resources/store/large-file.txt
${small file hash}      05a5b8a8b0280cf985e8de1f0cc1a980
${large file hash}      57e673156276e884bcb0207ee22e5a84
${large file hash 1}    ${small file hash}
${large file hash 2}    e44f0589587d1377cb68cf3166eb611e

** Keywords **
Create index from "${path}" and return lines
    ${result}=  Run process  ${indexer bin} --init\=true ${path}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}
    ${index lines}      Split to lines  ${result.stdout}
    [Return]    ${index lines}

Create index from "${path}" and save it to "${output path}"
    ${result}=  Run process  ${indexer bin} --init\=true ${path} | tee ${output path}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Store index "${index}" to "${store dir}" and return lines
    ${result}=  Run process  ${store bin} --maxBlockSize ${max block size} ${store dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${lines}=   Split to lines  ${result.stdout}
    [Return]   ${lines}

Store index "${index}" to "${store dir}" and save output to "${output path}"
    ${result}=  Run process  ${store bin} --maxBlockSize ${max block size} ${store dir} < ${index} | tee ${output path}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Restore index "${index}" from "${store dir}" to "${restore dir}"
    ${result}=  Run process  ${restore bin} ${store dir} ${restore dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Restore index dry run "${index}" from "${store dir}" to "${restore dir}"
    ${result}=  Run process  ${restore bin} -dry-run ${store dir} ${restore dir} < ${index}  shell=TRue
    Log many  ${result.stdout}
    Log many  ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

** Variables **
${differ bin}           bin/kopi-diff
${indexer bin}          bin/kopi-index
${store bin}            bin/kopi-store
${restore bin}          bin/kopi-restore

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

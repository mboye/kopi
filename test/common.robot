** Variables **
${differ bin}           bin/kopi-diff
${indexer bin}          bin/kopi-index
${store bin}            bin/kopi-store
${restore bin}          bin/kopi-restore
${manifest bin}         bin/kopi-manifest

${store dir}            ${TEMPDIR}/simple_store_data
${index}                ${TEMPDIR}/index
${max block size}       64

${backup source dir}    test/resources/store
${small file}           test/resources/store/small-file.txt
${large file}           test/resources/store/large-file.txt
${small file hash}      a2e94dfda3eb76fdb96649cea308ec07dde3243c
${large file hash}      aa2695f48e73a2e338dd20a6a501a8ad2d39c757
${large file hash 1}    ${small file hash}
${large file hash 2}    94794e6c56dca74fd44cc77e693523233c6022af
${empty file}           test/resources/store/empty-file
${empty file hash}      318f7856ed3b2050346de77f571469876e18e59f

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

Store index "${index}" with encryption to "${store dir}" and return lines
    ${result}=  Run process  ${store bin} --encrypt --maxBlockSize ${max block size} ${store dir} < ${index}  shell=True
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

Store index "${index}" with encryption to "${store dir}" and save output to "${output path}"
    ${result}=  Run process  ${store bin} --encrypt --maxBlockSize ${max block size} ${store dir} < ${index} | tee ${output path}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Restore index "${index}" from "${store dir}" to "${restore dir}"
    ${result}=  Run process  ${restore bin} ${store dir} ${restore dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Restore index "${index}" with encryption from "${store dir}" to "${restore dir}"
    ${result}=  Run process  ${restore bin} --decrypt ${store dir} ${restore dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

Restore index dry run "${index}" from "${store dir}" to "${restore dir}"
    ${result}=  Run process  ${restore bin} -dry-run ${store dir} ${restore dir} < ${index}  shell=TRue
    Log many  ${result.stdout}
    Log many  ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

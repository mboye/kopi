** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py
Resource    common.robot

Test Setup     Begin test
Test Teardown  End test

** Test Cases **
Write manifest
    Create index from "${backup source dir}" and save it to "${index}"

    ${result}=  Run process  ${manifest bin} write ${store dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${manifest id}=   Should match regexp  ${result.stderr}  (?m).*created manifest.*id=(.+)  groups=1
    File should exist           ${store dir}/manifests/${manifest id}

    ${data}=                    Get decompressed file   ${store dir}/manifests/${manifest id}
    ${lines}                    Split to lines          ${data}
    Length should be            ${lines}                5

    ${header}=                       Get from list   ${lines}  0
    Should be valid manifest header  ${header}

Write manifest with encryption
    Create index from "${backup source dir}" and save it to "${index}"

    ${result}=  Run process  ${manifest bin} write ${store dir} --encrypt < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${manifest id}=   Should match regexp  ${result.stderr}  (?m).*created manifest.*id=(.+)  groups=1
    File should exist           ${store dir}/manifests/${manifest id}

    Run keyword and expect error    *Not a gzipped file*
    ...     Get decompressed file   ${store dir}/manifests/${manifest id}

Read manifest with encryption
    Create index from "${backup source dir}" and save it to "${index}"

    ${description}=     Generate random string
    ${result}=          Run process  ${manifest bin} write ${store dir} --description\="${description}" --encrypt < ${index}  shell=True
    Log many            ${result.stdout}
    Log many            ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${manifest id}=   Should match regexp  ${result.stderr}  (?m).*created manifest.*id=(.+)  groups=1

    ${result}=  Run process  ${manifest bin} read ${store dir} ${manifest id} --decrypt  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${description read}=   Should match regexp  ${result.stderr}  .*description=(.+)${SPACE}  groups=1
    Should be equal as strings  ${description}  ${description read}

    ${lines}=           Split to lines  ${result.stdout}
    Length should be    ${lines}    4

Read manifest
    Create index from "${backup source dir}" and save it to "${index}"

    ${description}=     Generate random string
    ${result}=          Run process  ${manifest bin} write ${store dir} --description\="${description}" < ${index}  shell=True
    Log many            ${result.stdout}
    Log many            ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${manifest id}=   Should match regexp  ${result.stderr}  (?m).*created manifest.*id=(.+)  groups=1

    ${result}=  Run process  ${manifest bin} read ${store dir} ${manifest id}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${match}  ${description read}=   Should match regexp  ${result.stderr}  .*description=(.+)${SPACE}  groups=1
    Should be equal as strings  ${description}  ${description read}

    ${lines}=           Split to lines  ${result.stdout}
    Length should be    ${lines}    4

** Keywords **
Begin test
    Create directory        ${store dir}
    Copy file               test/resources/salt  ${store dir}/salt

End test
    Remove directory    ${store dir}  recursive=True

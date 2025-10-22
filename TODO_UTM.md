# WINDOWS11_ARM64_ISO_BUILDER



goal: create a bootable Windows 11 ARM64 ISO from Microsoft’s official UUP servers and register it in UTM.

NOTE needs validation !!! 

steps:
  - query: https://git.uupdump.net/uup-dump/json-api/listid.php?search=Windows+11&arch=arm64&sortByDate=1
    save: builds.json
    extract: response.builds[0].uuid -> BUILD_ID
  - query: https://git.uupdump.net/uup-dump/json-api/get.php?id=${BUILD_ID}&lang=en-us&edition=Core
    save: manifest.json
    extract: response.files -> FILES
  - for each FILE in FILES:
      download: FILE.url -> downloads/${BUILD_ID}/
      verify: sha1(FILE.path) == FILE.sha1
  - run:
      cmd: ./converter/convert.sh downloads/${BUILD_ID}/ output/Windows11_ARM64_${BUILD_ID}.iso
      check: file_exists(output/Windows11_ARM64_${BUILD_ID}.iso)
  - run:
      cmd: /Applications/UTM.app/Contents/MacOS/utmctl create --name Win11_ARM64_${BUILD_ID} --arch arm64 --memory 4096 --disk 64G --iso output/Windows11_ARM64_${BUILD_ID}.iso
  - run:
      cmd: /Applications/UTM.app/Contents/MacOS/utmctl start Win11_ARM64_${BUILD_ID}

output:
  iso: output/Windows11_ARM64_${BUILD_ID}.iso
  vm:  Win11_ARM64_${BUILD_ID}

validation:
  - all FILE.url domains end with mp.microsoft.com
  - sha1 verified
  - iso size > 5GB
  - iso boots in UTM

license:
  user must activate Windows legally
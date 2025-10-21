# todo

### UTM

In infra repo, we need to add UTM to the pkg/dep, so its alwas installed. 


### utm windows

From your Mac, we need to then use UTM to install windows 11, with the exact config we need. 

We need to do all this from golang.

https://github.com/naveenrajm7/
https://github.com/naveenrajm7/packer-plugin-utm 

THis will then allow use to git clone our infra project, and then run our code there for testing.


### utm Winget

When we are inside the windows machine that UTM creates, we need to call winget to install the same base things that we have on a mac

- git
- go
- vscoee
- vscode extensions.

that will allow use to run all the code on our infra repo.

We need to do all this from golang. 

github.com/crazy-max/winget-pkg

https://github.com/mbarbita/go-winget/blob/main/go-winget.go


### Winget MDM

https://github.com/jantari/rewinged

Might be useful later ... It Servers winget Manifeests.

SO i could host this, pop my package manifests in, and then install off the Server on any Windows machine ?
Must run behind caddy https.





# Mod Releaser

## Potential feature

 * git tag -> <component>/<version>
 * git push tags
 * find dagger modules + get only modified
 * implement automatic bump version


input:
 component name
 git repo
call:
 minor/patch/major

list git tag to search the latest release of a component
then bump the version if not found create a 0.1.0
push the git tag
publish with dagger
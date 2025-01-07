@ECHO OFF
:: This file is to create Stash Docker containers, and should be copied and called from the following path:
:: C:\Users\MyUserName\AppData\Local\Docker\wsl\CreateContainer.cmd
:: Example usage: 
::				CreateContainer.cmd MyContainerName "stashapp/stash:latest" 9998
:: Example with shared mount paths: 
::				CreateContainer.cmd ContainerName1 "stashapp/stash:latest" 9991 C:\MySharedMountPath C:\Another\Shared\Folder
:: Example adding Stash IMAGE and container:
::				CreateContainer.cmd NewContainer27.2 "stashapp/stash:v0.27.2" 9997 IMAGE
::		Note: The image name (stashapp/stash:v0.27.2) must be an image name listed in following link: https://hub.docker.com/r/stashapp/stash/tags
:: Example with DLNA:
::				CreateContainer.cmd MyDLNA272 "stashapp/stash:v0.27.2" 9996 C:\downloads DLNA
:: Example skipping docker-compose:
::					CreateContainer.cmd ContainerName "stashapp/stash:v0.26.2" 9992 C:\Videos SKIP
set NewContainerName=%1
:: Example image arguments:stashapp/stash:latest, stashapp/stash:v0.27.2, stashapp/stash:v0.26.2
set Image=%2
:: Example Port Numbers: 9999, 9990, 9991, 9995, 9998
set STASH_PORT=%3
:: The SharedMountPath's variables are optional arguments, and can be empty.
:: Use SharedMountPath's to specify shared paths that are mounted as READ-ONLY.
:: Example SharedMountPath: C:\Videos, E:\MyMedia, Z:\MyVideoCollections, C:\Users\MyUserName\Videos, C:\Users\MyUserName\download
set SharedMountPath=%4
set SharedMountPath2=%5
set SharedMountPath3=%6
set SharedMountPath4=%7
set SharedMountPath5=%8
set SharedMountPath6=%9
shift
shift
shift
shift
shift
set SharedMountPath7=%5
set SharedMountPath8=%6
set SharedMountPath9=%7
set SharedMountPath10=%8
set SharedMountPath11=%9
set SkipDockerCompose=
set DLNAFunctionality="no"
set PullDockerStashImage=
set MountAccess=:ro
if /I [%SharedMountPath%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath=)
if /I [%SharedMountPath%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath=)
if /I [%SharedMountPath%]==[IMAGE]	 (set PullDockerStashImage=yes) & (set SharedMountPath=)
if /I [%SharedMountPath%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath=)
if /I [%SharedMountPath2%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath2=)
if /I [%SharedMountPath2%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath2=)
if /I [%SharedMountPath2%]==[IMAGE]  (set PullDockerStashImage=yes) & (set SharedMountPath2=)
if /I [%SharedMountPath2%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath2=)
if /I [%SharedMountPath2%]==[WRITE]	 (set MountAccess=) & (set SharedMountPath2=)
if /I [%SharedMountPath3%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath3=)
if /I [%SharedMountPath3%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath3=)
if /I [%SharedMountPath3%]==[IMAGE]  (set PullDockerStashImage=yes) & (set SharedMountPath3=)
if /I [%SharedMountPath3%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath3=)
if /I [%SharedMountPath3%]==[WRITE]	 (set MountAccess=) & (set SharedMountPath3=)
if /I [%SharedMountPath4%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath4=)
if /I [%SharedMountPath4%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath4=)
if /I [%SharedMountPath4%]==[IMAGE]  (set PullDockerStashImage=yes) & (set SharedMountPath4=)
if /I [%SharedMountPath4%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath4=)
if /I [%SharedMountPath4%]==[WRITE]	 (set MountAccess=) & (set SharedMountPath4=)
if /I [%SharedMountPath5%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath5=)
if /I [%SharedMountPath5%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath5=)
if /I [%SharedMountPath5%]==[IMAGE]  (set PullDockerStashImage=yes) & (set SharedMountPath5=)
if /I [%SharedMountPath5%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath5=)
if /I [%SharedMountPath5%]==[WRITE]	 (set MountAccess=) & (set SharedMountPath5=)
if /I [%SharedMountPath6%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath6=)
if /I [%SharedMountPath6%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath6=)
if /I [%SharedMountPath6%]==[IMAGE]	 (set PullDockerStashImage=yes) & (set SharedMountPath6=)
if /I [%SharedMountPath6%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath6=)
if /I [%SharedMountPath6%]==[WRITE]  (set MountAccess=) & (set SharedMountPath6=)
if /I [%SharedMountPath7%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath7=)
if /I [%SharedMountPath7%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath7=)
if /I [%SharedMountPath7%]==[IMAGE]	 (set PullDockerStashImage=yes) & (set SharedMountPath7=)
if /I [%SharedMountPath7%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath7=)
if /I [%SharedMountPath7%]==[WRITE]  (set MountAccess=) & (set SharedMountPath7=)
if /I [%SharedMountPath8%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath8=)
if /I [%SharedMountPath8%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath8=)
if /I [%SharedMountPath8%]==[IMAGE]	 (set PullDockerStashImage=yes) & (set SharedMountPath8=)
if /I [%SharedMountPath8%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath8=)
if /I [%SharedMountPath8%]==[WRITE]  (set MountAccess=) & (set SharedMountPath8=)
if /I [%SharedMountPath9%]==[DLNA] 	 (set DLNAFunctionality=yes) & (set SharedMountPath9=)
if /I [%SharedMountPath9%]==[SKIP] 	 (set SkipDockerCompose=yes) & (set SharedMountPath9=)
if /I [%SharedMountPath9%]==[IMAGE]	 (set PullDockerStashImage=yes) & (set SharedMountPath9=)
if /I [%SharedMountPath9%]==[PULL] 	 (set PullDockerStashImage=yes) & (set SharedMountPath9=)
if /I [%SharedMountPath9%]==[WRITE]  (set MountAccess=) & (set SharedMountPath9=)
if /I [%SharedMountPath10%]==[DLNA]  (set DLNAFunctionality=yes) & (set SharedMountPath10=)
if /I [%SharedMountPath10%]==[SKIP]  (set SkipDockerCompose=yes) & (set SharedMountPath10=)
if /I [%SharedMountPath10%]==[IMAGE] (set PullDockerStashImage=yes) & (set SharedMountPath10=)
if /I [%SharedMountPath10%]==[PULL]  (set PullDockerStashImage=yes) & (set SharedMountPath10=)
if /I [%SharedMountPath10%]==[WRITE] (set MountAccess=) & (set SharedMountPath10=)
if /I [%SharedMountPath11%]==[DLNA]  (set DLNAFunctionality=yes) & (set SharedMountPath11=)
if /I [%SharedMountPath11%]==[SKIP]  (set SkipDockerCompose=yes) & (set SharedMountPath11=)
if /I [%SharedMountPath11%]==[IMAGE] (set PullDockerStashImage=yes) & (set SharedMountPath11=)
if /I [%SharedMountPath11%]==[PULL]  (set PullDockerStashImage=yes) & (set SharedMountPath11=)
if /I [%SharedMountPath11%]==[WRITE] (set MountAccess=) & (set SharedMountPath11=)


:: If user incorrectly enters below arguments instead of Stash-Port, fetch the values, and let CHECK_STASH_PORT get the required Stash-Port.
if /I [%STASH_PORT%]==[DLNA] 	(set DLNAFunctionality=yes) & (set STASH_PORT=)
if /I [%STASH_PORT%]==[SKIP] 	(set SkipDockerCompose=yes) & (set STASH_PORT=)
if /I [%STASH_PORT%]==[IMAGE]	(set PullDockerStashImage=yes) & (set STASH_PORT=)
if /I [%STASH_PORT%]==[PULL] 	(set PullDockerStashImage=yes) & (set STASH_PORT=)
echo SkipDockerCompose = %SkipDockerCompose% ; DLNAFunctionality = %DLNAFunctionality%
set DockerComposeFile="docker-compose.yml"

if [%NewContainerName%]==[] goto :MissingArgumentNewContainerName
goto :HaveVariableNewContainerName
:MissingArgumentNewContainerName
set /p NewContainerName="Enter the new container name: "
if [%NewContainerName%]==[] call:ExitWithError 160 "ERROR_BAD_ARGUMENTS"
:HaveVariableNewContainerName

if [%Image%]==[] goto :MissingArgumentImage
goto :HaveVariableImage
:MissingArgumentImage
set /p Image="Enter the image name: "
if [%Image%]==[] call:ExitWithError 160 "ERROR_BAD_ARGUMENTS"
:HaveVariableImage

:CHECK_STASH_PORT
if [%STASH_PORT%]==[] goto :MissingArgumentSTASH_PORT
IF 1%STASH_PORT% NEQ +1%STASH_PORT% goto STASH_PORT_NOT_NUMERIC
goto :HaveVariableSTASH_PORT
:STASH_PORT_NOT_NUMERIC
echo Error ******************
echo Argument #3 requires a numeric value for Stash-Port.  You entered "%STASH_PORT%" instead.  Please enter a numberic value for Stash Port.
:MissingArgumentSTASH_PORT
set STASH_PORT=
set /p STASH_PORT="Enter the Stash port number: "
if [%STASH_PORT%]==[] call:ExitWithError 160 "ERROR_BAD_ARGUMENTS"
goto :CHECK_STASH_PORT
:HaveVariableSTASH_PORT

if exist %NewContainerName%\ (
  echo %NewContainerName% already exists. 
) else (
  echo creating folder %NewContainerName%
  mkdir %NewContainerName%
)
cd %NewContainerName%
echo DockerComposeFile=%DockerComposeFile%; NewContainerName=%NewContainerName%; Image=%Image%; STASH_PORT=%STASH_PORT%; DLNAFunctionality=%DLNAFunctionality%; SharedMountPath=%SharedMountPath%; SharedMountPath1=%SharedMountPath1%; SharedMountPath2=%SharedMountPath2%
echo services:> %DockerComposeFile%
echo   stash:>> %DockerComposeFile%
echo     image: %Image%>> %DockerComposeFile%
echo     container_name: %NewContainerName%>> %DockerComposeFile%
echo     restart: unless-stopped>> %DockerComposeFile%
if [%DLNAFunctionality%]==[yes] goto :DoDLNA_Functionality
echo     ports:>> %DockerComposeFile%
echo       - "%STASH_PORT%:9999">> %DockerComposeFile%
goto :SkipDLNA_Functionality
:DoDLNA_Functionality
echo     network_mode: host>> %DockerComposeFile%
:SkipDLNA_Functionality
echo     logging:>> %DockerComposeFile%
echo       driver: "json-file">> %DockerComposeFile%
echo       options:>> %DockerComposeFile%
echo         max-file: "10">> %DockerComposeFile%
echo         max-size: "2m">> %DockerComposeFile%
echo     environment:>> %DockerComposeFile%
echo       - STASH_STASH=/data/>> %DockerComposeFile%
echo       - STASH_GENERATED=/generated/>> %DockerComposeFile%
echo       - STASH_METADATA=/metadata/>> %DockerComposeFile%
echo       - STASH_CACHE=/cache/>> %DockerComposeFile%
if [%DLNAFunctionality%]==[yes] goto :DoDLNA_Functionality_pt2
echo       - STASH_PORT=9999>> %DockerComposeFile%
goto :SkipDLNA_Functionality_pt2
:DoDLNA_Functionality_pt2
echo       - STASH_PORT=%STASH_PORT%>> %DockerComposeFile%
:SkipDLNA_Functionality_pt2
echo     volumes:>> %DockerComposeFile%
echo       - /etc/localtime:/etc/localtime:ro>> %DockerComposeFile%
echo       - ./config:/root/.stash>> %DockerComposeFile%
echo       - ./data:/data>> %DockerComposeFile%
echo       - ./metadata:/metadata>> %DockerComposeFile%
echo       - ./cache:/cache>> %DockerComposeFile%
echo       - ./blobs:/blobs>> %DockerComposeFile%
echo       - ./generated:/generated>> %DockerComposeFile%
if [%SharedMountPath%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath%:/external%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath2%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath2%:/external2%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath3%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath3%:/external3%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath4%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath4%:/external4%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath5%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath5%:/external5%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath6%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath6%:/external6%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath7%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath7%:/external7%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath8%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath8%:/external8%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath9%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath9%:/external9%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath10%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath10%:/external10%MountAccess%>> %DockerComposeFile%
if [%SharedMountPath11%]==[] goto :SkipSharedMountPaths
echo       - %SharedMountPath11%:/external11%MountAccess%>> %DockerComposeFile%
:SkipSharedMountPaths

if [%SkipDockerCompose%] NEQ [] goto :DoNot_DockerCompose
if [%PullDockerStashImage%] NEQ [yes] goto :SkipPullDockerStashImage
docker pull %Image%
:SkipPullDockerStashImage
docker-compose up -d
:DoNot_DockerCompose
cd ..
Goto :eof

:ExitWithError
Echo Error: Exiting with error code %1 and error message %2
Exit %~1
Goto :eof

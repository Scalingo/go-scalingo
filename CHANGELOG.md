# Changelog

## v2.4.4

* Display request ID in the debug logs

## v2.4.3

* Add SSH to Region to configure SSH endpoint for SSH-based operations

## v2.4.2

* UsersSelf is now based on AuthenticationService not API anymore
* Add RegionsList() to get the list of available platform regions

## v2.4.1

* Implement new methods on Client: `EventTypesList` and `EventCategoriesList`

## v2.4.0

* Update `Notifier` methods to match the API: ie. accept SelectedEventIDs as input and output

## v2.3.0

* Update `NotifierUpdate`/`NotifierProvision` params, Add missing field `Notifier.SendAllAlerts`

## v2.2.1

* Add `Description` and `LogoURL` to `NotificationPlatform` struct

## v2.2.0

* Add `OperationsShowFromURL(url string) (*Operation, error)` to ease the use
  of the Operation URL returned after a Scale/Restart action
* Add `OperationStatus` and `OperationType` types with the right constants
* Remove `Plan.TextDescription` which was never used

## v2.1.1

* Add StaticTokenGenerator in ClientConfig to ensure retrocompatibility

## v2.1.0

* StacksList() to list available runtime stacks
* Add AppsSetStack() to update the stack of an app

## v2.0.0

* Integration of Database API authentication
* Ability to query backup/logs of addon
* Add missing Addon#Status field

## v1.5.2

* Remove os.Exit(), reliquat from split between CLI and client.
* Update wording
* Fix display of alert table

## v1.5.1

* Update deps

## v1.5.0

* Add AppsForceHTTPS
* Add AppsStickySession
* Add AppID in App subresources
* Collaborator.Status is now of type CollaboratorStatus, and constants are defined

## v1.4.1

* Add UserID to Collaborator

## v1.4.0

* Add Fullname for `User` model
* Ability to create an email notifier
* Access to one-off audit logs

## v1.3.2

* Add events NewAlert, Alert, DeleteAlert, NewAutoscaler, DeleteAutoscaler

## v1.3.0

* Change keys endpoint to point to the authentication service instead of the main API
* Add `GithubLinkService` implementation

## v1.2.0

* Refactoring, use interface instead of private struct

## v1.1.0

* API Token methods + authentication

## v1.0.0

* Initial tag

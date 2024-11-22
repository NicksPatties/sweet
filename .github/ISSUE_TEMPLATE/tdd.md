---
name: Design document
about: Clearly define the scope of a large feature
title: '[TDD] '
labels: ['design document']
assignees: ''
---

# _Feature Title_ TDD

Author: _Your Name_
Date: _YYYY-MM-DD_ (updated), _YYYY-MM-DD_ (created)

## Overview

_Short description of the change and/or feature that is being proposed._

## Key User Stories and tasks

- _User story one title_
  - User story description
    - *As a ______, I need _____, so that ______* 
  - Priority
    - (must have/should have/could have)
  - Tasks needed to achieve goal (aka: "User Flow")
    1. _Task one_
      - _Links to details about task one_
    2. _Task two_
      - _Links to details about task two_
      - _Another link..._
    3. _Task three_

## What

<details>

### Background

_Describe any historical context that would be needed to understand the document, including legacy considerations._

### Terminology

If the document uses any special words or terms, list them here.

### Non-Goals

- _Non-goal one_
- _Non-goal two_
- _Non-goal three_

### Technical requirements

#### New Entrypoints

_"Entrypoints" could be web server endpoints, commands in a CLI application, or something smiliar._

- _Entrypoint name_
  - _Request type_
    - (GET/POST/...)
  - Why is this new entrypoint needed?
    - _Reason for why this entrypoint is necessary_ 
  - Description of input/output contract
    - _If request X is sent in with Y, then Z should be returned_
    - _If command X is executed with Y flag, then Z should be returned_

#### Additions/Changes to Existing Entrypoints

- _Entrypoint name_
  - _Request type_
    - (GET/POST/...)
  - Why is this change required?
    - _Reason for why this entrypoint is necessary_ 
  - Description of input/output contract
    - _If request X is sent in with Y, then Z should be returned_
    - _If command X is executed with Y flag, then Z should be returned_

#### Storage Model Additions/Changes

_These are changes to ways that data is stored in the application._

- _Model name_
  - Additions and/or changes
    - _Change one_
    - _Change two_
    - _Change three_
  - How will existing data be handled?
    - _This is how existing data will be handled_

#### UI Screens/Components

- _UI component title_
  - _A description of the UI component_
  - Mockups: _Place mockups of UI component here_

#### Data Handling and Privacy

- _Type of data_
  - _A description of the data_
  - Why do we need to store this data?
    - _A reason why_
  - Anonymized?
    - _Yes/No_
  - Can the user opt out?
    - _Yes/No_
  - Wipeout (user delete) policy
    - _What should happen to the data when the user deletes their account_
  - Takeout (data export) policy
    - _What should happen when the user exports the data._

#### Connection to Existing Work

_Is there prior work that can be drawn on for this feature? If so, link it here._

#### Other Requirements

_Add any other technical considerations that have not been mentioned yet._

### Testing Plan

- Test Description: _What should the test demonstrate?_
  - Initial setup steps
    1. _Do a thing_
    2. _Do another thing_
  - Test steps
    1. _Click this thing_
      - [ ] _The expected behavior should be this_
    2. _Type this thing_
      - [ ] _The expected behavior should be that_

### Remaining Open Questions

_Add any questions that you do not konw the answer to at this moment. If there are no more questions, you can leave this section blank._

- _Question 1_
- _Question 2_ 

> **STOP.
> Review the [WHAT section](#what) again. If the proposal looks good, then move on.**

<details>
<summary>Review considerations</summary>

- Should we tackle this problem at this time? (yes/no)
- Do the requirements in this TDD match those in the [user stories](#key-user-stories-and-tasks)? (yes/no)
- Can the project be done using only existing patterns in the codebase? (yes/no)
  - Are there any approaches/insights that might be useful for designing the solution? 
- Are all assertions properly justified (with links to sources/proof, if appropriate)? (yes/no)
- Does the testing plan validate all required user stories in the product spec? (yes/no)
</details>

</details>

## How

<details>

### Existing Status Quo

_If a similar task was performed previously, then link it here._

### Solution Overview

_Add a sketch of a sequence or dependency diagram here. Give a high level view of the possible solution._

#### Third party libraries

- [_libraryname:version_](link)
  - Why is it needed?
    - _It's needed to do XYZ_
  - License
    - _MIT, Apache 2.0, etc._

#### Third party services

_Things like external APIs (Google Maps, OpenAI, etc.) are stated here._

- [_Service name_](link)
  - Why is it needed?
    - _Reason_
  - What is the plan if the dependency fails?
    - _Mitigation strategy_

### Architectual Decisions

#### Decision 1: *Decision name*

The following options have been considered.

- _Option one_
- _Option two_
- _..._

Among these, _this one_ may be the best option, because:
- _Reason one_
- _Reason two_
- _..._

The options are compared in detail in the table below

|    |  Option one | Option two |
| --- | --- | --- |
| _Consideration one_ | _Option one impact_ | _Option two impact_ |
| _Consideration two_ | _Option one impact_ | _Option two impact_ |
| ... | | |

#### Decision 2: *Decision name*

_Repeat as many decision sections as needed._

### Risks and mitigations

- _Describe the potential risk here_
  - Mitigation: _This is the mitigation for the above risk_

> **STOP. 
> Review from the top of the [HOW section](#how).**

<details>
<summary>Review considerations</summary>

- Is the proposed solution understandable to others? (yes/no)
- Is the architectural decision analysis solid, and are the conclusions well-reasoned and supported by the data? (yes/no)
  - Are there other key architectural decisions that should be considered (but haven’t been)? (yes/no)
- Will the proposed approach scale? (yes/no)
- Are there any potential red flags or risks (particularly around security and compatibility) in the proposed solution that need further investigation? (yes/no)
- Are the architectural patterns being implemented consistent with the rest of the affected codebase? (yes/no)
  - If new patterns are being introduced, do they set the right precedents and are they well-designed? (yes/no)
</details>

### Implementation approach

_Describe the changes that you intend to make to the codebase. **You are not required to write code at this point.** Pseudocode is encouraged._

_Also, describe explicit contracts that should be followed in the implementation, including function definitions and comments that describe what they should be doing._

#### Storage Model

_Describe how the data is stored. Some questions below to consider:_

_How their IDs are generated?_
_What are their fields (including descriptions)?_
_Any constraints on validation?_
_Is the data a source of truth, or is it derived from something else?_
_How can the user export this data? (Takeout policy)_
_How is the data removed if a user no longer wants to use the app? (Wipeout policy)_
_How will you query this data? Will you require pagination?_
_What are the method signatures that will be needed? Include doc comments, but you don't need to include code._

#### Storage Model Migrations 

_Describe any required migrations for your data storage. Be especially cognizant if you require backfilling any rows with existing data!_

#### Domain Objects

_Describe any new objects related to your business logic that will be created. Give its name, fields (descriptions and types), and methods that will be supported. Include doc comments, but no need to include code._

#### User flow

_For each [user story and task](#key-user-stories-and-tasks), create a flow diagram for any interaction with entrypoints into the application. Consider the following:_

- _if making a web application/feature, consider the series of request/response operations (endpoint, handlers, pseudocode describing what happens on the server, etc.)_
- _domain controllers_
- _utiliy functions_
- _core business logic_

#### UI Changes

_For each [user story and task](#key-user-stories-and-tasks), create a series of screens that the user will see when using the application. Take extra consideration to ensure that the following are considered:_

- _screen constraints (height, width, landscape, vertical orientation)_
- _light and dark themes_
- _high pixel density devices_
- _left-to-right and right-to-left languages_
- _accessibility concerns_
- _how data is bound or transformed to show these UI components_

#### Documentation Changes

_List changes to documentation for the feature to be clearly described. Be sure to include changes to READMEs, help messages, wiki pages, and so on._ 

### Metrics plan

_What metrics will be recorded? How do they supplement the test cases, or verify that the behavior is as expected?_

</details>

## Implementation Plan

_Break down the list of tasks into issues here._

## Launch Plan

_Share how this feature will be released to the public. Will this feature be behind a feature flag? Will it be released in an upcoming update?_

### Future work

_Is there any work that needs to be done in the future to maintain this feature, or is related to this feature in some way?_


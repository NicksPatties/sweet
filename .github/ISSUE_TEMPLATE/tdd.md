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

## Problem Statement

*Use "as a _____, I would like _____" format.*

## What

<details>

### Background

_Describe any historical context that would be needed to understand the document, including legacy considerations._

### Terminology

If the document uses any special words or terms, list them here.

### Non-Goals

_Describe a list of items that are **not** included in this design document_.

- _Item one_
- _Item two_
- _Item three_

### Technical requirements

_The below sections contains some sample considerations in your design._

#### Additions/Changes to Endpoints

_These could be web server endpoints, or commands in a CLI application._

#### Storage Model Additions/Changes

_These are changes to ways that data is stored in the application._

- _Model name_
  - Description of changes
    - _These are the things that will be changed_
  - How will existing data be handled?
    - _This is how existing data will be handled_

#### UI Screens/Components

- _UI component title_
  - Description:
    - _A description of the UI component_

_Place mockups of UI component here_

#### Data Handling and Privacy

- _Type of data_
  - Description
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
    _1. Do a thing_
    _2. Do another thing_
  - Test steps
    _1. Click this thing_
      - [ ] _The expected behavior should be this_
    _2. Type this thing_
      - [ ] _The expected behavior should be that_

### Remaining Open Questions

_Add any questions that you do not konw the answer to at this moment. If there are no more questions, you can leave this section blank._

- _Question 1_
- _Question 2_ 

**Review the What section and verify that it looks good before moving to the next section.**

</details>

## How

**Only do this step once you have a solid understanding of what needs to be accomplished in the [What section](#what).**

<details>

### Existing Status Quo

_If a similar task was performed previously, then state it here._

### Solution Overview

_Add a sketch of a sequence or dependency diagram here. Give a high level view of the possible solution._

#### Third party libraries

- _Library name_
  - [link](_link to dependency_)
  - Why is it needed?
    - _It's needed to do XYZ_
  - License
    - _MIT, Apache 2.0, etc._

#### Third party services

_Things like external APIs (Google Maps, OpenAI, etc.) are stated here._

- _Service name_
  - [link](_link to dependency_)
  - Why is it needed?
    - _Reason_
  - What is the plan if the dependency fails?
    - _Mitigation strategy_

### Architectual Decisions

#### Decision 1: *Decision name*

The following options have been considered.

*Option 1*
*Option 2*

_Now, list the impacted topics, and compare the options together._

- (Example) Data storage efficiency
  - ~Option 1~
    - _some data_
  - Option 2
    - _some other data_

_Mark options that are no longer in consideration with a ~strikethrough~._

#### Decision 2: *Decision name*

_Repeat as many decision sections as needed._

### Risk and mitigations

- Potential risk: _there is something bad that could happen here._
  - Mitigation: _This is the mitigation for the above risk_

**At this point, you should verify your data and risks have properly been assessed before moving to implementation**.

### Implementation approach

_Describe the changes that you intend to make to the codebase. **You are not required to write code at this point.** Pseudocode is encouraged._

_Also, describe explicit contracts that should be followed in the implementation, including function definitions and comments that describe what they should be doing._

#### Data Model

_Describe how the data is stored. This could include a description of the database schema._

#### Interface/API Definitions

_Describe how the various components talk to each other. For example, if there are REST endpoints, describe the endpoint URL and the format of the data and parameters used._

#### Business Logic

_If the design requires any non-trivial algorithms or logic, describe them._

#### Migration Strategy

_If the design incurs non-backwards-compatible changes to an existing system, describe the process whereby entities that depend on the system are going to migrate to the new design._

### Metrics plan

_What metrics will be recorded? How do they supplement the test cases, or verify that the behavior is as expected?_

</details>

## Implementation Plan

_Break down the list of tasks into issues here._

## Launch Plan

_Share how this feature will be released to the public. Will this feature be behind a feature flag? Will it be released in an upcoming update?

### Future work

_Is there any work that needs to be done in the future to maintain this feature, or is related to this feature in some way? This is similar to the [non-goals section](#non-goals) of the document.


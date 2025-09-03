# JIRA Ticket

https://n-able.atlassian.net/browse/

# Change Description

A high-level description of the changes made (including refactoring)

# Risk

- Category: Low/Medium/High
- Reason: Why the risk category was chosen
- Mitigation: For Medium or High, how will the risk be mitigated?

# Developer Testing

Upon completion the Developer should be confident they have done all the testing they can before handing over to QA.

## Unit Testing

- A high-level description of unit test coverage of the change.
- Identify areas not covered by unit testing to be covered by developer integration testing.

## Local Integration Testing

- A grid or bullet list of scenarios to be integration tested by the Developer to prove the change.
- Identify areas not covered by integration testing to be covered by QA.
- Identify any short-cuts taken that might not reflect actual system usage.
- **NOT** a replacement for integration testing by QA - an addition to find obvious bugs.

## Testing Still Required
- Handover to QA
- Make it clear to QA what you expect them to test.
- Include areas where regression coverage may be required.
- This is not the Developer defining the QA test plan - we are contributing ideas to it.
- QA ultimately still define the test cases for the change.

---

[//]: # 'Links should be alphabetised for readability'
[conventional commit]: https://www.conventionalcommits.org/en/v1.0.0/
[signed]: https://portal.n-able.dev/docs/northstar/system/northstar/GitHub/Commit-Signing/

# Rollback that does not fully repair the incident

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

Possible causes include:

- The code change also modified shared configuration.
- New code wrote data the old version cannot interpret.
- An external side effect already occurred.
- Some instances never completed rollback.
- A cache retains incompatible data.
- A migration changed a dependency used by services outside the initial target
  set.
- The apparent correlation with the deployment was coincidental.

The follow-up "Would you roll back?" should lead to an assessment of
reversibility and evidence, not an automatic answer.

Explain what exactly would be restored, what remains changed, whether the old
version has enough capacity, and how to validate recovery across affected
cohorts. A successful rollback API response is not proof that customer-visible
behavior recovered.

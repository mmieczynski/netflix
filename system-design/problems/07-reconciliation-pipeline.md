# Reconciliation pipeline

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

For reconciliation, prepare:

> System A says 10,000 items should exist downstream. System B reports 9,997.
> How do you identify the missing items, determine whether they are delayed or
> lost, and repair them without duplicating effects?

That question connects naturally to the public payment-pipeline mock. Be ready
to discuss stable identities, expected-versus-observed state, age thresholds,
authoritative records, safe replay, and an audit trail for repair.

# Distributed counting

[Preparation index](../README.md) · [Interview evidence](../interview-context.md)

Practice scenario for Netflix L5 Developer Platform. The follow-ups are
preparation questions, not confirmed Netflix interview questions.

For counting, start with a simple race:

> The cap is five. Two servers each observe a count of four and both approve
> one more operation.

An atomic increment alone does not fix a separate non-atomic check. Explain
where the admission decision and reservation become indivisible.

Then prepare questions about:

- Fixed versus rolling windows.
- Duplicate and delayed events.
- Reservations for operations that may never complete.
- Regional partitions.
- Whether a strict global bound is required.
- The availability or utilization cost of a single authority or preallocated
  regional budgets.

Netflix's counter abstraction distinguishes different durability and
consistency requirements. Use it to study why an eventually accurate count and
a strict admission limit are different requirements.
[Netflix's Distributed Counter Abstraction](https://netflixtechblog.com/netflixs-distributed-counter-abstraction-8d0c45eb66b2)

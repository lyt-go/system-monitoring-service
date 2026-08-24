# Bug Reproduction

## Problem

A staged sample batch retains the caller's slice and sample pointers across the later commit.

## Trigger

Stage two samples, reuse the input slice and mutate its first sample, then commit and read the stored batch twice.

## Error

The committed IDs and values follow the reused input, and mutating a list result changes the next query.

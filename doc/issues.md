#### Frontend
- [ ] Reservation management
- [ ] Queue Prompt
- [x] Event status view

#### Backend
- [ ] Email confirmations
- [ ] Email reminders
- [ ] Configuration of queueing and reservation rules
- [ ] Scheduling
- [ ] Check what data is actually needed for modifying and cancelling reservations
- [ ] Check that all access to common data is thread safe. Use a consistent strategy for when (on what level) locks are used
- [ ] DB encryption
    - [ ] password in ENV variable

### Issues
- [ ] **Reject any empty request**
- [x] **Error handling for frontend API**
- [ ] Database culling trigger?
    - [ ] Clients by Epoch accessed:
    - [ ] Delete clients whose last accessed is over 2 years (configurable)
- [ ] Slow requests issue (Hash algorythm?)
- [ ] Session key handling requires a redesign to incorporate fingerprint, hashing and performant searching
- [ ] Light weight event list should probably be cached and not re-created for each request
    - map out everything the list events do
- [ ] Should there be a limit to reservations in queue?
    - Now you can only make a reservation, if there is at least one free spot in the timeslot
- [x] Delete event doesn't update instantly
- [ ] If timeslots are given out of order, the wrongly sorted item is silently lost

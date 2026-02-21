#### Frontend
- [x] Event browse view
- [x] Reservation form
- [ ] Reservation management
- [ ] Queue Prompt
- [x] Account management
- [x] Secure session implementation
- [x] Admin view
- [x] Event creation
- [x] Event modification
- [ ] Event status view
- [x] Show Draft status on event (admin)

#### Backend
- [x] Secure session implementation
- [x] Making reservations
- [x] Editing reservations
- [x] Crating events
- [ ] Email confirmations
- [x] Queueing
- [ ] Email reminders
- [ ] Configuration of queueing and reservation rules
- [ ] Scheduling
- [x] Make sure reservations cannot be modified or removed, unless user matches the reservation
- [ ] Check what data is actually needed for modifying and cancelling reservations
- [ ] Check that all access to common data is thread safe. Use a consistent strategy for when (on what level) locks are used

### Issues
- [ ] **Reject any empty request**
- [ ] **Error handling for frontend API**
- [ ] Database culling trigger?
    - [ ] Clients by Epoch accessed:
    - [ ] Delete clients whose last accessed is over 2 years (configurable)
- [ ] Slow requests issue (Hash algorythm?)
- [ ] Session key handling requires a redesign to incorporate fingerprint, hashing and performant searching
- [ ] Light weight event list should probably be cached and not re-created for each request
    - map out everything the list events do
- [ ] Should there be a limit to reservations in queue?
    - Now you can only make a reservation, if there is at least one free spot in the timeslot

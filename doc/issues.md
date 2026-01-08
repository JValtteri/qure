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
- [ ] Show Draft status on event (admin)

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

### Issues

- [x] **Making an anonymous reservation with existing email gives a session cookie for existing acount!**
- [x] **Session cookie remains valid after logout**
    - [x] Backend
    - [x] Frontend
- [ ] **Reject any empty request**
- [ ] **Error handling for frontend API**
- [x] ID is used as password, but cannot handle non-unique password
- [x] Cannot create a custom password
- [x] Difference between ID and Key?
- [ ] Database culling trigger?
    - [ ] Clients by Epoch accessed:
    - [ ] Delete clients whose last accessed is over 2 years (configurable)
- [x] Reservation Test has lots of disabled tests
- [ ] ID sometimes contains an unhashed password

- [x] Make variables configurable
    - [x] Password length
    - [x] Username length
    - [x] Time to keep temp account
    - [x] Time to keep reservation
        - [x] Tentative
        - [x] Past start
- [x] Password length not checked on changed password

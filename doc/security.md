# QuRe Reservation System Security Documentation

## Index

- [Collected Data](#collected-data)
- [Fingerprinting](#fingerprinting)
- [Access Conditions](#access-conditions)
- [Potential Issues](#potential-issues)


## Collected Data

|      Collected      | Reason                                                                                                                                                                                                                                                                                            |      Stored      |
| :-----------------: | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------: |
|  email (username)   | for reservations, for delivering reservation vahvistus                                                                                                                                                                                                                                            |    plain text    |
|      password       | this should be unique and not a problem.                                                                                                                                                                                                                                                          |      hashed      |
| browser fingerprint | For added security: for the duration of a session (max 30 days) is used as a second factor to validate the session cookie used to resume sessions.<br>As the fingerprint is stored hashed, the information isn't visible to anyone and reverse engineering it would be difficult for an attacker. |      hashed      |
|         IP          | Only in case of error IP may be present in server logs. Data is only in RAM and is purger when server is restarted. This is a feature of the networking library, not our code                                                                                                                     | server log (RAM) |
|     SessionKey      | For securing active user session. The key doesn't contain any identifiable information, other than that it is unique. <br>It alone won't give access to an account. A matching fingerprint is required.\*                                                                                           |    plain text    |
|                     |                                                                                                                                                                                                                                                                                                   |                  |

*\*) This is true only if `EXTRA_STRICT_SESSIONS` are enabled, or if sessions are redesigned to improve fingerprinting performance*

## Fingerprinting

Backend collects standard information sent by the web browser client and forms a quasi unique fingerprint of the client. The fingerprint is hashed, so that the individual data elements can no longer be inspected. A fingerprint is generated whenever a user would resume their session (i.e. do actions while logged in). The fingerprint is compared to stored fingerprint to make sure the client is the same as the one with whom the session was formed with. This allows the server to detect if someone is trying to use someone else's session key. This is passive fingerprinting.

In active fingerprinting, JavaScript is sent to the client and run there to collect much more information about the capabilities and features of the client system. This would be more accurate, but also more invasive. It is used in the wild, but QuRe does not use active fingerprinting.

See [Cover Your Tracks by Electronic Frontier Foundation](https://coveryourtracks.eff.org/) for more information on browser fingerprinting.

|                    |   bits   |
| ------------------ | :------: |
| UserAgent          |   4.5    |
| Accepted           |   10.5   |
| Language           |   6.5    |
|         **Total:** | **21.5** |

With 21.5 bits of identifying information, you have about **2.9 million** possible combinations, though some are more common than others. It means that it's not a foolproof way to identify a single person. The data can be spoofed and it is easy to collect by any other site a user visits, but in combination with a cryptographically random session cookie, the session should be reasonably safe from brute force attacks.

## Access Conditions

| Action                 | Required authentication                       |
| ---------------------- | --------------------------------------------- |
| Login user             | **credentials** / **reservation ID**          |
| View reservations      | **credentials** OR associated **session key** |
| Modify reservation     | associated **reservation ID** only            |
| Cancel reservation     | associated **reservation ID** only            |
| View user reservations | ***                                           |
| View all reservations  | ***                                           |
| Create event           | ***                                           |
| Edit event             | ***                                           |
| Modify user            | ***                                           |
| Delete user            | ***                                           |

## Potential Issues

- [ ] **Session keys in searchable clientsBySession map are unhashed**
- [ ] ReservationID is sometimes used as a password
    - [ ] Reservation IDs are stored unhashed
- [ ] emails are stored unencrypted
- [ ] With `EXTRA_STRICT_SESSIONS` turned off, fingerprinting is ineffective against session forgery
- [ ] **Session cookie parameters need to be adjusted to achieve maximum security**
    - [ ] `Secure` cookie parameter
    - [ ] others [Cookie security](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#cookies)
- [ ] **Frontend has very little input validation**
    - [ ] Dates
    - [ ] Times
    - [ ] Reservation sizes
- [ ] No request throttling mechanism implemented. Backend is suceptible to DoS attacks.
- [ ] Number of sessions is unlimited (Risk: out of memory crash)
- [ ] Number of users is unlimited (Risk: out of memory crash)
- [ ] Request performance is low (Risk: DoS)

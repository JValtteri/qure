import './PrivacyPolicy.css'

function PrivacyPolicy() {
    return (
        <div className="tooltip">privacy policy
            <div className="tooltiptext">
                By clicking accept, you give us permission to:
                <p>
                    1. Create an account for you to save you reservations.
                </p>
                <p>
                    2. Place a cookie on your system to authenticate your session
                </p>
                <p>
                    3. To store your browser fingerprint as a way to validate the
                    authenticity of your session
                </p>
                <p>
                    Following GDPR you have a right to revoke concent:
                </p>
                <p>
                    You can at any time log out, to remove the session cookie and
                    the accompanying fingerprint information.
                </p>
                <p>
                    You can at any time delete your account to remove your
                    registration  information, reservations and remaining session.
                </p>
                <p>
                    You can view the information stored about you on your user
                    settings page
                </p>
                <p>
                    In case of errors, your IP address may be collected as part of
                    monitoring the health of the service. The information is stored
                    and used only to diagnose any issues in the service or abuse.
                </p>
                <p>
                    Your data is not given to third parties unless required by law.
                </p>
            </div>
        </div>
    )
}

export default PrivacyPolicy;

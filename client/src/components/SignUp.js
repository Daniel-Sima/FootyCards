import React, { useState } from "react"
// import theme from "./Theme"
import LoginGif from "./../assets/login.gif"
import Sheet from '@mui/joy/Sheet/Sheet';
// import SportsSoccerIcon from '@mui/icons-material/SportsSoccer';
import Typography from '@mui/joy/Typography';
// import Divider from '@mui/joy/Divider'
// import Avatar from '@mui/joy/Avatar'
import Box from '@mui/joy/Box';
import Button from '@mui/joy/Button';
// import IconButton from '@mui/joy/IconButton';
// import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import FormControl from '@mui/joy/FormControl';
import FormLabel from '@mui/joy/FormLabel';
import Input from '@mui/joy/Input';
import { Link } from "react-router-dom";
import Alert from "./Alert"

const SignUp = () => {
    const [pseudo, setPseudo] = useState("");
    const [clubName, setClubName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [alertMessage, setAlertMessage] = useState("");
    const [alertClassName, setAlertClassName] = useState("d-none");
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = (event) => {
        event.preventDefault();
        
        // Cas particulier 
        if ((clubName === "") || (clubName === "") || (email === "") || (password === "")) {
            setAlertClassName("alert-danger");
            setAlertMessage("Missing fields!");
            event.target.reset();
            return
        }

        setIsSubmitting(true);

        // Build the request payload
        let payload = {
            pseudo: pseudo,
            clubName: clubName,
            email: email,
            password: password,
        }

        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(payload),
        }

        fetch(`/register`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setAlertClassName("alert-danger");
                    setAlertMessage(data.message);
                } else {
                    setEmail("");
                    setPassword("");
                    setPseudo("");
                    setClubName("");
                    setAlertClassName("alert-success");
                    setAlertMessage("Successful registration, you can log in now!");
                    event.target.reset();
                }
            })
            .catch(err => {
                setAlertClassName("alert-danger");
                setAlertMessage(err);
            })
            .finally(() => {
                setIsSubmitting(false);
              })
    }

    return (
        <>
            <Box
                sx={{
                    backgroundImage: `linear-gradient(0deg, #334ec1 0%, #000a19 40%)`,
                    alignItems: 'center',
                    justifyContent: 'center',
                    display: 'flex',
                    height: '100vh',
                }}
            >
                {
                    <>
                        <Sheet
                            variant="outlined"
                            sx={{
                                width: 700,
                                py: 3,
                                px: 2,
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 2,
                                borderRadius: 'sm',
                            }}
                        >
                            <Box
                                sx={{
                                    position: 'absolute',
                                    transform: 'translate(80%, -125%)',
                                    left: '-450px',
                                    width: alertMessage === "" ? '100%' : '105%',
                                    height: alertMessage === "" ? '30%' : '20%',
                                    backgroundImage: `url(${LoginGif})`,
                                    backgroundSize: 'contain',
                                    backgroundRepeat: 'no-repeat',

                                }}
                            />

                            <div>
                                <Typography align="center" level="h2" component="h1">
                                    Create account
                                </Typography>
                                {/* <Typography align="center" level="h5">Create account</Typography> */}
                            </div>

                            <form onSubmit={handleSubmit}>
                                <FormControl>
                                    <FormLabel sx={{ fontSize: '15px' }}>Pseudo</FormLabel>
                                    <Input
                                        name="pseudo"
                                        type="text"
                                        onChange={(event) => setPseudo(event.target.value)}
                                    />
                                </FormControl>

                                <FormControl>
                                    <FormLabel sx={{ pt: 3, fontSize: '15px' }}>Club name</FormLabel>
                                    <Input
                                        name="clubName"
                                        type="text"
                                        onChange={(event) => setClubName(event.target.value)}
                                    />
                                </FormControl>

                                <FormControl>
                                    <FormLabel sx={{ pt: 3, fontSize: '15px' }}>Email</FormLabel>
                                    <Input
                                        name="email"
                                        type="email"
                                        onChange={(event) => setEmail(event.target.value)}
                                    />
                                </FormControl>

                                <FormControl>
                                    <FormLabel sx={{ pt: 3, fontSize: '15px' }}>Password</FormLabel>
                                    <Input
                                        name="password"
                                        type="password"
                                        onChange={(event) => setPassword(event.target.value)}
                                    />
                                </FormControl>

                                <div style={{ display: 'flex', justifyContent: 'center' }}>
                                    <Button type='submit' disabled={isSubmitting} sx={{ mt: 3 }} style={{ width: 200, height: 50, fontSize: '18px' }}  >
                                        {isSubmitting ? "In progress..." : "Create account"}
                                    </Button>
                                </div>

                                <Typography
                                    endDecorator={<Link to="/login">Log in</Link>}
                                    level="body1"
                                    style={{ display: 'flex', justifyContent: 'center' }}
                                >
                                    Have an account?
                                </Typography>

                            </form>
                            <Alert
                                message={alertMessage}
                                className={alertClassName}
                            />
                        </Sheet>
                    </>
                }
            </Box>
        </>
    )
}

export default SignUp;
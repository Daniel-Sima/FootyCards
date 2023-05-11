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
import { Link, useNavigate } from "react-router-dom";
import Alert from "./Alert"
// import MyContext from "./MyContext";
// import Header from "./Header/Header";

const Login = (props) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [alertMessage, setAlertMessage] = useState("");
    const [alertClassName, setAlertClassName] = useState("d-none");
    
    const navigate = useNavigate();

    const handleSubmit = (event) => {
        event.preventDefault();

        // Build the request payload
        let payload = {
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

        fetch(`https://footycards-production-39e4.up.railway.app/authenticate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setAlertClassName("alert-danger");
                    setAlertMessage(data.message);
                    event.target.reset();
                } else {
                    props.setJwtToken(data.access_token); // 
                    localStorage.setItem('jwtToken', data.access_token); 
                    localStorage.setItem('email', email);  
                    setEmail("");
                    setPassword("");
                    setAlertClassName("d-none");
                    setAlertMessage("");
                    props.toggleRefresh(true);
                    navigate("/team");
                }
            })
            .catch(err => {
                setAlertClassName("alert-danger");
                setAlertMessage(err);
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
                                    left: '-500px',
                                    width:  '100%',
                                    height: alertMessage === "" ? '45%' : '40%',
                                    backgroundImage: `url(${LoginGif})`,
                                    backgroundSize: 'contain',
                                    backgroundRepeat: 'no-repeat',

                                }}
                            />

                            <div>
                                <Typography align="center" level="h2" component="h1">
                                    Welcome!
                                </Typography>
                                <Typography align="center" level="body1">Sign in to continue.</Typography>
                            </div>

                            <form onSubmit={handleSubmit}>
                                <FormControl>
                                    <FormLabel sx={{ fontSize: '15px' }}>Email</FormLabel>
                                    <Input
                                        name="email"
                                        type="email"
                                        // placeholder="adminBG@email.com"
                                        onChange={(event) => setEmail(event.target.value)}
                                    />
                                </FormControl>

                                <FormControl>
                                    <FormLabel sx={{ pt: 3, fontSize: '15px' }}>Password</FormLabel>
                                    <Input
                                        name="password"
                                        type="password"
                                        // placeholder="1234"
                                        onChange={(event) => setPassword(event.target.value)}
                                    />
                                </FormControl>

                                <div style={{ display: 'flex', justifyContent: 'center' }}>
                                    <Button sx={{ mt: 3 }} style={{ width: 100, height: 50, fontSize: '18px' }} type='submit'>
                                        Log in
                                    </Button>
                                </div>

                                <Typography
                                    endDecorator={<Link to="/sign-up">Sign up</Link>}
                                    level="body1"
                                    style={{ display: 'flex', justifyContent: 'center' }}
                                >
                                    Don't have an account?
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

export default Login;
import React from 'react';
import FormControl from '@mui/joy/FormControl';
import Radio from '@mui/joy/Radio';
import RadioGroup from '@mui/joy/RadioGroup';
import Sheet from '@mui/joy/Sheet';
import Typography from '@mui/joy/Typography';
import Box from '@mui/joy/Box';

const RadioButton = (props) => {
    // const [club, setClub] = useState("")

    // useEffect(() => {
    //     props.handleEquipes(club)

    // }, [club]);

    return (
        <FormControl style={{ fontSize: '24px', position: 'fixed', top: '125px', right: '10px', fontFamily: 'Impact, sans-serif', fontWeight: 'bold' }}  >
            <RadioGroup
                overlay
                name="equipes"
                orientation="vertical"
                sx={{ mt: 1 }}
            >
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(3, 1fr)',
                        gap: 1,
                    }}
                >

                    {props.availableClubs && Array.from(props.availableClubs).map((nom) => (
                        <Sheet
                            component="label"
                            key={nom}
                            variant="outlined"
                            sx={{
                                p: 2,
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'start',
                                boxShadow: 'sm',
                                borderRadius: 'md',
                                bgcolor: '#96d2ff',
                                borderColor: 'blue',
                                gap: 1.5,
                                width: '100px',
                                height: '100px',
                            }}
                        >
                            <Radio
                                onClick={() => props.handleEquipes(nom)} 
                                variant="outlined"
                                value={`person${nom}`}
                                sx={{
                                    mt: -1,
                                    mr: -1,
                                    mb: -1.5,
                                    alignSelf: 'flex-end',
                                    '--Radio-actionRadius': (theme) => theme.vars.radius.md,
                                }}
                            />
                            <Typography sx={{ ml: -1, color: 'black' }} level="body1">{nom}</Typography>
                        </Sheet>
                    ))}
                </Box>
            </RadioGroup>
        </FormControl>
    );
}

export default RadioButton;
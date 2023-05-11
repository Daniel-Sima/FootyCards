import {extendTheme} from '@mui/joy/styles';


const theme = extendTheme({
    components: {
      JoyCard: {
        styleOverrides: {
          root: ({ ownerState, theme }) => ({
            ...(ownerState.color === 'primary' && {
              // color: "#0A99F7",
              background: "#0A99F7",
              backgroundColor: "linear-gradient(#e66465, #9198e5)",
            }),
          }),
        },
      },
      JoySheet: {
        styleOverrides: {
          root: ({ ownerState, theme }) => ({
            ...(ownerState.color === 'primary' && {
              color: theme.vars.palette.text.secondary,
              backgroundColor: "#00022e",
            }),
            ...(ownerState.color === 'secondary' && {
              color: theme.vars.palette.text.secondary,
              backgroundColor: "#00035b",
            }),
          }),
        },
      },
      JoyButton: {
        styleOverrides: {
          root: ({ ownerState, theme }) => ({
            ...(ownerState.color === 'primary' && {
              color: "#FFFFFF",
              backgroundColor: "#00022e",
            }),
          }),
        },
      },
    },
  });

export default theme;
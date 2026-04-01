import { alpha, createTheme, responsiveFontSizes } from "@mui/material/styles";

export const wbPalette = {
  berry: "#CB11AB",
  orchid: "#B42371",
  plum: "#990099",
  mulberry: "#481173",
  midnight: "#14061F",
  mist: "#F7EFFF",
  success: "#5CD6A4",
  warning: "#F9B24E",
  danger: "#FF6F91",
  paper: "rgba(31, 16, 46, 0.62)",
  paperStrong: "rgba(18, 8, 28, 0.84)",
  border: "rgba(255, 255, 255, 0.14)",
};

let theme = createTheme({
  cssVariables: true,
  palette: {
    mode: "dark",
    primary: {
      main: wbPalette.berry,
      light: "#E45AD0",
      dark: wbPalette.plum,
    },
    secondary: {
      main: wbPalette.mulberry,
      light: "#7B40A6",
      dark: "#24083B",
    },
    success: {
      main: wbPalette.success,
    },
    warning: {
      main: wbPalette.warning,
    },
    error: {
      main: wbPalette.danger,
    },
    background: {
      default: "#080510",
      paper: wbPalette.paper,
    },
    text: {
      primary: "#FFF7FF",
      secondary: alpha(wbPalette.mist, 0.72),
    },
    divider: alpha("#FFFFFF", 0.08),
  },
  shape: {
    borderRadius: 4,
  },
  typography: {
    fontFamily: '"Manrope", "Segoe UI", sans-serif',
    h1: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
      letterSpacing: "-0.05em",
    },
    h2: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
      letterSpacing: "-0.04em",
    },
    h3: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
      letterSpacing: "-0.03em",
    },
    h4: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
      letterSpacing: "-0.02em",
    },
    h5: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
    },
    h6: {
      fontFamily: '"Space Grotesk", "Manrope", sans-serif',
      fontWeight: 700,
    },
    button: {
      fontWeight: 700,
      textTransform: "none",
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: "#080510",
          backgroundImage: [
            "radial-gradient(circle at 8% 10%, rgba(203, 17, 171, 0.28), transparent 28%)",
            "radial-gradient(circle at 88% 16%, rgba(72, 17, 115, 0.34), transparent 24%)",
            "radial-gradient(circle at 56% 88%, rgba(180, 35, 113, 0.20), transparent 32%)",
            "linear-gradient(135deg, #07030C 0%, #12041C 42%, #1D0A2D 100%)",
          ].join(","),
          color: "#FFF7FF",
        },
        "::selection": {
          backgroundColor: alpha(wbPalette.berry, 0.42),
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backdropFilter: "blur(22px)",
          backgroundImage: "none",
          backgroundColor: wbPalette.paper,
          border: `1px solid ${wbPalette.border}`,
          boxShadow: "0 18px 60px rgba(0, 0, 0, 0.24)",
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: "rgba(16, 7, 24, 0.82)",
          borderRight: `1px solid ${alpha("#FFFFFF", 0.08)}`,
          backdropFilter: "blur(26px)",
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: alpha("#12061D", 0.68),
          backdropFilter: "blur(18px)",
          borderBottom: `1px solid ${alpha("#FFFFFF", 0.08)}`,
          boxShadow: "none",
        },
      },
    },
    MuiButton: {
      defaultProps: {
        disableElevation: true,
      },
      styleOverrides: {
        root: {
          borderRadius: 18,
          paddingInline: 18,
          minHeight: 44,
        },
        containedPrimary: {
          background: "linear-gradient(135deg, #CB11AB 0%, #8B168A 100%)",
          boxShadow: "0 12px 28px rgba(203, 17, 171, 0.28)",
        },
        outlined: {
          borderColor: alpha("#FFFFFF", 0.14),
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 999,
          border: `1px solid ${alpha("#FFFFFF", 0.08)}`,
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        variant: "outlined",
        fullWidth: true,
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 18,
          backgroundColor: alpha("#FFFFFF", 0.03),
        },
        notchedOutline: {
          borderColor: alpha("#FFFFFF", 0.12),
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        root: {
          borderBottom: `1px solid ${alpha("#FFFFFF", 0.08)}`,
        },
        head: {
          color: alpha(wbPalette.mist, 0.76),
          fontWeight: 700,
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          backgroundColor: alpha("#15081F", 0.88),
          backgroundImage: "none",
        },
      },
    },
  },
});

theme = responsiveFontSizes(theme, {
  factor: 2.4,
});

export const appTheme = theme;

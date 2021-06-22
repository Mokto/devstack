import { extendTheme, ThemeConfig } from "@chakra-ui/react"
// 2. Add your color mode config
const config : ThemeConfig = {
  useSystemColorMode: false,
}

const theme = extendTheme({
  config,
  styles: {
    global: {
      body: {
        bg: "gray.600",
        color: "white",
      },
    },
  },
})
export default theme

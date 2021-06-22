import { useColorMode, Button } from '@chakra-ui/react';
import { memo } from 'react';

export const IndexPage = memo(() => {
  const { colorMode, toggleColorMode } = useColorMode()
    return (
        <div className="App">
          <header className="App-header">
            {/* <img src={logo} className="App-logo" alt="logo" /> */}
            <p>
              Edit <code>src/App.tsx</code> and save to reload.
            </p>
            <a
              className="App-link"
              href="https://reactjs.org"
              target="_blank"
              rel="noopener noreferrer"
            >
              Learn React
            </a>
            <Button onClick={toggleColorMode}>
              Toggle {colorMode === "light" ? "Dark" : "Light"}
            </Button>
          </header>
        </div>
      );
})

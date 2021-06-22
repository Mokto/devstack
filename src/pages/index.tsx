import { useColorMode, Button } from '@chakra-ui/react';
import { useReducer } from 'react';
import { useEffect } from 'react';
import { memo } from 'react';
import { websocket } from '../services/websocket';
import {
  default as AnsiUp
} from 'ansi_up';

const ansi_up = new AnsiUp();

export const IndexPage = memo(() => {
  const { colorMode, toggleColorMode } = useColorMode();
  const [logs, dispatchNewLog] = useReducer((state: string[], action: string) => {
    return [...state, action]
  }, []);

  useEffect(() => {
    const onLog = (data: any) => {
      dispatchNewLog(data.message);
      console.log(data.service)
    }
    websocket.on('log', onLog)
    return () => websocket.off('log', onLog)
  }, [])

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

            <hr />
            {logs.map(log => {
              return (
                <div dangerouslySetInnerHTML={{__html: ansi_up.ansi_to_html(log)}}></div>
              )
            })}
          </header>
        </div>
      );
})

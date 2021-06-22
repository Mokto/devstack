import { useColorMode, Button, SimpleGrid, Box } from '@chakra-ui/react';
import { useCallback, useReducer } from 'react';
import { useEffect } from 'react';
import { memo } from 'react';
import { websocket } from '../services/websocket';
import {
  default as AnsiUp
} from 'ansi_up';
import { useState } from 'react';
import { useLayoutEffect } from 'react';

const ansi_up = new AnsiUp();
ansi_up.use_classes = true;

interface Service {
  name: string;
  color: string;
  command: string;
}
interface Config {
  services: Service[]
}

interface LogData {
  message: string;
  service: Service;
}

export const IndexPage = memo(() => {
  const { colorMode, toggleColorMode } = useColorMode();
  const [config, setConfig] = useState<Config>();
  const [logs, dispatchNewLog] = useReducer((state: LogData[], action: LogData) => {
    return [...state, action]
  }, []);

  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/config');
      const result = await data.json();
      setConfig(result);
    };

    getData();
    
  }, [])

  // logs
  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/logs');
      const result: LogData[] = await data.json();
      
      console.log(result.forEach((log) => dispatchNewLog(log)))
    };

    getData();

    const onLog = (data: any) => {
      dispatchNewLog(data);
    }
    websocket.on('log', onLog)
    return () => websocket.off('log', onLog)
  }, [])

  useLayoutEffect(() => {
    const logsElement = document.getElementById('logs');
    if (logsElement) {
      logsElement.scrollTop = logsElement?.scrollHeight;
    }
  }, [logs])

  const applyPrefix = useCallback((log: LogData) => {
    return (
      <>
      <span style={{display: 'inline-block', width: 100}} className={`ansi-${log.service.color}-fg`}>[{log.service.name}]</span>
        {' - '}
        <span dangerouslySetInnerHTML={{__html: ansi_up.ansi_to_html(log.message)}}></span>
      </>
    )
  }, [])

    return (
        <div className="App">
          <header className="App-header">
            <Button onClick={toggleColorMode}>
              Toggle {colorMode === "light" ? "Dark" : "Light"}
            </Button>

            <hr />
            {config?.services && (
              <SimpleGrid columns={4} spacing={10}>
                {config?.services.map(service => {
                  return (
                    <Box padding={10} border={"1px solid red"} key={service.name}>
                      {service.name}


                      <Button onClick={() => {
                        fetch(`http://localhost:9111/restart/${service.name}`, {method:'post'})
                      }}>
                        Restart
                      </Button>
                    </Box>
                  )
                })}
              </SimpleGrid>
            )}
                <Box padding={10} height={500} overflowY="scroll" id="logs">
                  {logs.map((log, index) => {
                    return (
                      <div key={index}>
                        {applyPrefix(log)}
                      </div>
                    )
                  })}
                </Box>
          </header>
        </div>
      );
})


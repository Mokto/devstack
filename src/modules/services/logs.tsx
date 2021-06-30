import React, {
  useCallback,
  useReducer,
  useEffect,
  memo,
  useLayoutEffect,
} from 'react';
import { default as AnsiUp } from 'ansi_up';
import { websocket } from '../../services/websocket';
import styled from 'styled-components';

const LogsContainer = styled.div`
  position: absolute;
  top: 0;
  bottom: 0;
  left: 300px;
  right: 0;
  overflow-y: scroll;
  padding: 10px 20px;
`;

const ansi_up = new AnsiUp();
ansi_up.use_classes = true;

interface Service {
  name: string;
  color: string;
  command: string;
  isWatching: boolean;
  isRunning: boolean;
  watchDirectories: string[];
}

interface LogData {
  message: string;
  service: Service;
}

export const Logs = memo(() => {
  const [logs, dispatchNewLog] = useReducer(
    (state: LogData[], action: LogData) => {
      return [...state, action];
    },
    []
  );

  // logs
  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/logs');
      const result: LogData[] = await data.json();

      result.forEach(log => dispatchNewLog(log));
    };

    getData();

    const onLog = (data: any) => {
      dispatchNewLog(data);
    };
    websocket.on('log', onLog);
    return () => websocket.off('log', onLog);
  }, []);

  useLayoutEffect(() => {
    const logsElement = document.getElementById('logs');
    if (logsElement) {
      logsElement.scrollTop = logsElement?.scrollHeight;
    }
  }, [logs]);

  const applyPrefix = useCallback((log: LogData) => {
    return (
      <>
        <span
          style={{ display: 'inline-block', width: 200 }}
          className={`ansi-${log.service.color}-fg`}
        >
          [{log.service.name}]
        </span>
        {' - '}
        <span
          dangerouslySetInnerHTML={{
            __html: ansi_up.ansi_to_html(log.message),
          }}
        ></span>
      </>
    );
  }, []);

  return (
    <LogsContainer id="logs">
      {logs.map((log, index) => {
        return <div key={index}>{applyPrefix(log)}</div>;
      })}
    </LogsContainer>
  );
});

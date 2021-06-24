import React, {
  useCallback,
  useReducer,
  useEffect,
  memo,
  useState,
  useLayoutEffect,
} from 'react';
import { websocket } from '../services/websocket';
import { default as AnsiUp } from 'ansi_up';

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
interface State {
  services: Service[];
}

interface LogData {
  message: string;
  service: Service;
}

export const IndexPage = memo(() => {
  const [state, setState] = useState<State>();
  const [logs, dispatchNewLog] = useReducer(
    (state: LogData[], action: LogData) => {
      return [...state, action];
    },
    []
  );

  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/state');
      const result = await data.json();
      setState(result);
      console.log(result);
    };

    getData();
  }, []);

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
    <div>
      {state?.services && (
        <div className="services-container">
          {state?.services.map(service => {
            return (
              <div
                key={service.name}
                className={`service-box ${
                  service.isRunning
                    ? 'service-is-running'
                    : 'service-is-stopped'
                } ${
                  service.watchDirectories && !service.isWatching
                    ? 'service-disabled-watching'
                    : ''
                }`}
              >
                <div className="title">{service.name}</div>
                {service.isRunning && (
                  <button
                    onClick={() => {
                      fetch(`http://localhost:9111/restart/${service.name}`, {
                        method: 'post',
                      });
                    }}
                  >
                    Restart
                  </button>
                )}
                <button
                  onClick={() => {
                    service.isRunning = !service.isRunning;
                    service.isWatching = service.isRunning;
                    setState({ ...state });
                    fetch(`http://localhost:9111/setRunning/${service.name}`, {
                      method: 'post',
                      body: JSON.stringify({
                        isRunning: service.isRunning,
                      }),
                      headers: new Headers({
                        'Content-Type': 'application/json',
                      }),
                    });
                  }}
                >
                  {service.isRunning ? 'Stop' : 'Start'} service
                </button>
                {service.watchDirectories && service.isRunning && (
                  <button
                    onClick={() => {
                      service.isWatching = !service.isWatching;
                      setState({ ...state });
                      fetch(
                        `http://localhost:9111/setWatching/${service.name}`,
                        {
                          method: 'post',
                          body: JSON.stringify({
                            isWatching: service.isWatching,
                          }),
                          headers: new Headers({
                            'Content-Type': 'application/json',
                          }),
                        }
                      );
                    }}
                  >
                    {service.isWatching ? 'Stop' : 'Start'} watching
                  </button>
                )}
              </div>
            );
          })}
        </div>
      )}
      <div className="logs" id="logs">
        {logs.map((log, index) => {
          return <div key={index}>{applyPrefix(log)}</div>;
        })}
      </div>
    </div>
  );
});

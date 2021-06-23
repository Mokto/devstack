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
}
interface Config {
  services: Service[];
}

interface LogData {
  message: string;
  service: Service;
}

export const IndexPage = memo(() => {
  const [config, setConfig] = useState<Config>();
  const [logs, dispatchNewLog] = useReducer(
    (state: LogData[], action: LogData) => {
      return [...state, action];
    },
    []
  );

  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/config');
      const result = await data.json();
      setConfig(result);
    };

    getData();
  }, []);

  // logs
  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/logs');
      const result: LogData[] = await data.json();

      console.log(result.forEach(log => dispatchNewLog(log)));
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
          style={{ display: 'inline-block', width: 100 }}
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
      {config?.services && (
        <div className="services-container">
          {config?.services.map(service => {
            return (
              <div key={service.name} className="service-box">
                <div className="title">{service.name}</div>
                <button
                  onClick={() => {
                    fetch(`http://localhost:9111/restart/${service.name}`, {
                      method: 'post',
                    });
                  }}
                >
                  Restart
                </button>
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

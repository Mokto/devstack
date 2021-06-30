import React, { useEffect, memo, useState } from 'react';
import styled, { css } from 'styled-components';
import { websocket } from '../../services/websocket';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faRedo,
  faStop,
  faPlay,
  faEye,
  faEyeSlash,
} from '@fortawesome/free-solid-svg-icons';

const StateContainer = styled.div`
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 300px;
  overflow-y: auto;
  padding: 10px 5px;
`;

const Service = styled.div<{ active: boolean; disabledWatching: boolean }>`
  height: 60px;
  border: 1px solid #aaa;
  border-right: 30px solid #ff5555;
  margin-bottom: 5px;
  padding-left: 10px;

  button {
    display: inline-block;
    padding: 0.25em 0.5em;
    border: 0.1em solid #ffffff;
    margin: 0 0.3em 0.3em 0;
    border-radius: 0.12em;
    box-sizing: border-box;
    text-decoration: none;
    font-family: 'Roboto', sans-serif;
    font-weight: 300;
    color: #ffffff;
    text-align: center;
    transition: all 0.2s;
    background: transparent;
    cursor: pointer;
  }

  button:hover {
    color: #000000;
    background-color: #ffffff;
  }

  .title {
    height: 30px;
    line-height: 30px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  ${({ active }) =>
    active &&
    css`
      border-right-color: #50fa7b;
    `}

  ${({ disabledWatching }) =>
    disabledWatching &&
    css`
      border-right-color: #ff79c6;
    `}
`;

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

export const State = memo(() => {
  const [state, setState] = useState<State>();

  useEffect(() => {
    const getData = async () => {
      const data = await fetch('http://localhost:9111/state');
      const result = await data.json();
      setState(result);
      console.log(result);
    };

    getData();
  }, []);

  // state
  useEffect(() => {
    const onChangeState = ({ isRunning, serviceName }: any) => {
      if (state) {
        const service = state.services.find(s => s.name === serviceName);
        if (service) {
          service.isRunning = isRunning;
        }
      }

      console.log(serviceName, isRunning);

      setState({ ...state } as any);
    };
    websocket.on('isRunning', onChangeState);
    return () => websocket.off('isRunning', onChangeState);
  }, [state]);

  if (!state?.services) {
    return null;
  }

  return (
    <StateContainer>
      {state?.services.map(service => {
        return (
          //   <div
          //     key={service.name}
          //     className={`service-box ${
          //       service.isRunning ? 'service-is-running' : 'service-is-stopped'
          //     } ${
          //       service.watchDirectories && !service.isWatching
          //         ? 'service-disabled-watching'
          //         : ''
          //     }`}
          //   >
          <Service
            key={service.name}
            active={service.isRunning}
            disabledWatching={
              service.watchDirectories &&
              service.isRunning &&
              !service.isWatching
            }
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
                <FontAwesomeIcon icon={faRedo} />
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
              {service.isRunning ? (
                <FontAwesomeIcon icon={faStop} />
              ) : (
                <FontAwesomeIcon icon={faPlay} />
              )}
            </button>
            {service.watchDirectories && service.isRunning && (
              <button
                onClick={() => {
                  service.isWatching = !service.isWatching;
                  setState({ ...state });
                  fetch(`http://localhost:9111/setWatching/${service.name}`, {
                    method: 'post',
                    body: JSON.stringify({
                      isWatching: service.isWatching,
                    }),
                    headers: new Headers({
                      'Content-Type': 'application/json',
                    }),
                  });
                }}
              >
                {service.isWatching ? (
                  <FontAwesomeIcon icon={faEyeSlash} />
                ) : (
                  <FontAwesomeIcon icon={faEye} />
                )}
              </button>
            )}
          </Service>
        );
      })}
    </StateContainer>
  );
});

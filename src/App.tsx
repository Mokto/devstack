import React, { memo, useState, useEffect } from 'react';
import { IndexPage } from './pages/index';
import { websocket } from './services/websocket';

export default memo(() => {
  const [isWebsocketConnected, setIsWebsocketConnected] = useState<boolean>(
    websocket.isConnected
  );
  useEffect(() => {
    const onChangeState = (isConnected: boolean) => {
      setIsWebsocketConnected(isConnected);
    };
    websocket.onConnectionEvent(onChangeState);

    return () => websocket.offConnectionEvent(onChangeState);
  }, []);

  if (!isWebsocketConnected) {
    return <h1>Not connected to server. Make sure it's online.</h1>;
  }

  return <IndexPage />;
});

import { Theme } from '@carbon/react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import MainLayout from './components/layout/MainLayout';
import OverviewPage from './pages/OverviewPage';
import TopicsPage from './pages/TopicsPage';
import ConsumerGroupsPage from './pages/ConsumerGroupsPage';
import BrokersPage from './pages/BrokersPage';
import './App.scss';

function App() {
  return (
    <Router>
      <Theme theme="g100">
        <MainLayout>
          <Routes>
            <Route path="/" element={<Navigate to="/overview" replace />} />
            <Route path="/overview" element={<OverviewPage />} />
            <Route path="/topics" element={<TopicsPage />} />
            <Route path="/consumer-groups" element={<ConsumerGroupsPage />} />
            <Route path="/brokers" element={<BrokersPage />} />
          </Routes>
        </MainLayout>
      </Theme>
    </Router>
  );
}

export default App;

#!/usr/bin/env python3
"""
Integration tests for MCP Server
Tests the server as a whole system
"""

import pytest
import asyncio
import json
from unittest.mock import AsyncMock, patch, MagicMock
from io import StringIO

from mcp.server.stdio import stdio_server
from mcp.server.models import InitializationOptions
from mcp.types import JSONRPCMessage, JSONRPCRequest, JSONRPCResponse

from server import server, SERVER_VERSION, main


class TestServerIntegration:
    """Integration tests for the complete server"""
    
    @pytest.mark.asyncio
    async def test_server_initialization(self):
        """Test server initialization with proper options"""
        init_options = InitializationOptions(
            server_name="my-mcp-server",
            server_version=SERVER_VERSION,
            capabilities=server.get_capabilities(
                notification_options=None,
                experimental_capabilities=None,
            ),
        )
        
        assert init_options.server_name == "my-mcp-server"
        assert init_options.server_version == SERVER_VERSION
        assert init_options.capabilities is not None
    
    @pytest.mark.asyncio
    async def test_server_capabilities(self):
        """Test that server reports correct capabilities"""
        capabilities = server.get_capabilities(
            notification_options=None,
            experimental_capabilities=None,
        )
        
        # Check that tools capability is present
        assert hasattr(capabilities, 'tools')
        # Add more capability checks as the server grows
    
    @pytest.mark.integration
    @pytest.mark.asyncio
    async def test_full_server_lifecycle(self):
        """Test complete server lifecycle simulation"""
        # This test simulates a full client-server interaction
        
        # Mock stdio streams
        mock_read_stream = AsyncMock()
        mock_write_stream = AsyncMock()
        
        # Create a mock request for listing tools
        list_tools_request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/list",
            "params": {}
        }
        
        # Mock the read stream to return our request
        mock_read_stream.read.return_value = json.dumps(list_tools_request).encode()
        
        # This is a simplified test - in a real scenario, you'd need to
        # properly mock the entire JSONRPC protocol flow
        
        # For now, just test that the server can be instantiated
        assert server is not None
        assert hasattr(server, 'list_tools')
        assert hasattr(server, 'call_tool')


class TestServerProtocol:
    """Test MCP protocol compliance"""
    
    @pytest.mark.asyncio
    async def test_tools_list_protocol(self):
        """Test that tools/list follows MCP protocol"""
        tools = await server._list_tools_handler()
        
        # Each tool should have required MCP fields
        for tool in tools:
            assert hasattr(tool, 'name')
            assert hasattr(tool, 'description')
            assert hasattr(tool, 'inputSchema')
            
            # Name should be a string
            assert isinstance(tool.name, str)
            assert len(tool.name) > 0
            
            # Description should be a string
            assert isinstance(tool.description, str)
            assert len(tool.description) > 0
            
            # Input schema should be a valid JSON schema
            assert isinstance(tool.inputSchema, dict)
            assert 'type' in tool.inputSchema
    
    @pytest.mark.asyncio
    async def test_tools_call_protocol(self):
        """Test that tools/call follows MCP protocol"""
        # Test valid tool call
        result = await server._call_tool_handler("get_server_version", {})
        
        # Result should be a list of content
        assert isinstance(result, list)
        assert len(result) > 0
        
        # Each content item should have proper structure
        for content in result:
            assert hasattr(content, 'type')
            assert hasattr(content, 'text')
            assert content.type == "text"
            assert isinstance(content.text, str)


class TestServerPerformance:
    """Performance tests for the server"""
    
    @pytest.mark.asyncio
    async def test_concurrent_tool_calls(self):
        """Test server handles concurrent tool calls"""
        # Create multiple concurrent calls
        tasks = []
        for i in range(10):
            task = server._call_tool_handler("get_server_version", {})
            tasks.append(task)
        
        # Execute all tasks concurrently
        results = await asyncio.gather(*tasks)
        
        # All should succeed
        assert len(results) == 10
        for result in results:
            assert len(result) == 1
            assert SERVER_VERSION in result[0].text
    
    @pytest.mark.asyncio
    async def test_rapid_tool_listing(self):
        """Test rapid tool listing calls"""
        # Make multiple rapid calls to list tools
        tasks = []
        for i in range(5):
            task = server._list_tools_handler()
            tasks.append(task)
        
        results = await asyncio.gather(*tasks)
        
        # All should return the same tools
        assert len(results) == 5
        for result in results:
            assert len(result) >= 1
            # All results should be identical
            assert result == results[0]


if __name__ == "__main__":
    # Run integration tests
    pytest.main([__file__, "-v", "-m", "integration"])
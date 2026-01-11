#!/usr/bin/env python3
"""
Test cases for MCP Server
"""

import pytest
import asyncio
from unittest.mock import AsyncMock, patch
from mcp.types import TextContent, Tool

# Import the server components
from server import (
    server,
    get_server_version,
    example_tool_handler,
    handle_list_tools,
    handle_call_tool,
    SERVER_VERSION
)


class TestServerVersion:
    """Test cases for server version functionality"""
    
    @pytest.mark.asyncio
    async def test_get_server_version(self):
        """Test getting server version"""
        result = await get_server_version({})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert result[0].type == "text"
        assert SERVER_VERSION in result[0].text
        assert "MCP Server Version:" in result[0].text
    
    @pytest.mark.asyncio
    async def test_get_server_version_with_arguments(self):
        """Test getting server version with arguments (should ignore them)"""
        result = await get_server_version({"unused_param": "value"})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert SERVER_VERSION in result[0].text


class TestToolListing:
    """Test cases for tool listing functionality"""
    
    @pytest.mark.asyncio
    async def test_handle_list_tools(self):
        """Test listing available tools"""
        tools = await handle_list_tools()
        
        assert isinstance(tools, list)
        assert len(tools) >= 1
        
        # Check that get_server_version tool is present
        version_tool = next((tool for tool in tools if tool.name == "get_server_version"), None)
        assert version_tool is not None
        assert isinstance(version_tool, Tool)
        assert version_tool.description == "Get the current version of this MCP server"
        assert version_tool.inputSchema["type"] == "object"
        assert version_tool.inputSchema["required"] == []
    
    @pytest.mark.asyncio
    async def test_tool_schema_structure(self):
        """Test that all tools have proper schema structure"""
        tools = await handle_list_tools()
        
        for tool in tools:
            assert hasattr(tool, 'name')
            assert hasattr(tool, 'description')
            assert hasattr(tool, 'inputSchema')
            assert isinstance(tool.name, str)
            assert isinstance(tool.description, str)
            assert isinstance(tool.inputSchema, dict)
            assert 'type' in tool.inputSchema
            assert 'properties' in tool.inputSchema
            assert 'required' in tool.inputSchema


class TestToolCalling:
    """Test cases for tool calling functionality"""
    
    @pytest.mark.asyncio
    async def test_handle_call_tool_version(self):
        """Test calling get_server_version tool"""
        result = await handle_call_tool("get_server_version", {})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert SERVER_VERSION in result[0].text
    
    @pytest.mark.asyncio
    async def test_handle_call_tool_version_with_args(self):
        """Test calling get_server_version tool with arguments"""
        result = await handle_call_tool("get_server_version", {"param": "value"})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert SERVER_VERSION in result[0].text
    
    @pytest.mark.asyncio
    async def test_handle_call_tool_version_none_args(self):
        """Test calling get_server_version tool with None arguments"""
        result = await handle_call_tool("get_server_version", None)
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert SERVER_VERSION in result[0].text
    
    @pytest.mark.asyncio
    async def test_handle_call_tool_unknown(self):
        """Test calling unknown tool raises ValueError"""
        with pytest.raises(ValueError, match="Unknown tool: nonexistent_tool"):
            await handle_call_tool("nonexistent_tool", {})
    
    @pytest.mark.asyncio
    async def test_handle_call_tool_empty_name(self):
        """Test calling tool with empty name raises ValueError"""
        with pytest.raises(ValueError, match="Unknown tool: "):
            await handle_call_tool("", {})


class TestExampleToolHandler:
    """Test cases for example tool handler"""
    
    @pytest.mark.asyncio
    async def test_example_tool_handler_with_param(self):
        """Test example tool handler with parameter"""
        result = await example_tool_handler({"param": "test_value"})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert result[0].type == "text"
        assert "Processed parameter: test_value" in result[0].text
    
    @pytest.mark.asyncio
    async def test_example_tool_handler_without_param(self):
        """Test example tool handler without parameter"""
        result = await example_tool_handler({})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert "Processed parameter: " in result[0].text
    
    @pytest.mark.asyncio
    async def test_example_tool_handler_empty_param(self):
        """Test example tool handler with empty parameter"""
        result = await example_tool_handler({"param": ""})
        
        assert len(result) == 1
        assert isinstance(result[0], TextContent)
        assert "Processed parameter: " in result[0].text


class TestServerIntegration:
    """Integration tests for the server"""
    
    def test_server_version_constant(self):
        """Test that SERVER_VERSION is properly defined"""
        assert isinstance(SERVER_VERSION, str)
        assert len(SERVER_VERSION) > 0
        assert "." in SERVER_VERSION  # Should be in format like "1.0.0"
    
    @pytest.mark.asyncio
    async def test_server_capabilities(self):
        """Test server capabilities"""
        capabilities = server.get_capabilities(
            notification_options=None,
            experimental_capabilities=None,
        )
        
        assert capabilities is not None
        # Add more specific capability tests as needed
    
    @pytest.mark.asyncio
    async def test_full_workflow(self):
        """Test complete workflow: list tools -> call tool"""
        # First, list tools
        tools = await handle_list_tools()
        assert len(tools) >= 1
        
        # Find the version tool
        version_tool = next((tool for tool in tools if tool.name == "get_server_version"), None)
        assert version_tool is not None
        
        # Call the version tool
        result = await handle_call_tool("get_server_version", {})
        assert len(result) == 1
        assert SERVER_VERSION in result[0].text


class TestErrorHandling:
    """Test cases for error handling"""
    
    @pytest.mark.asyncio
    async def test_malformed_tool_call(self):
        """Test handling of malformed tool calls"""
        with pytest.raises(ValueError):
            await handle_call_tool("invalid_tool_name", {})
    
    @pytest.mark.asyncio
    async def test_none_tool_name(self):
        """Test handling of None tool name"""
        with pytest.raises((ValueError, TypeError)):
            await handle_call_tool(None, {})


if __name__ == "__main__":
    # Run tests
    pytest.main([__file__, "-v"])